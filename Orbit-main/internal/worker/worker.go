package worker

import (
	"Orbit/internal/db"
	"Orbit/internal/rediskeys"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductPayload struct {
	ID        bson.ObjectID `bson:"_id"`
	EventID   bson.ObjectID `bson:"eventId"`
	SellerID  bson.ObjectID `bson:"sellerId"`
	Title     string        `bson:"title"`
	Price     float64       `bson:"price"`
	Frequency int           `bson:"frequency"`
	Currency  string        `bson:"currency"`
	EndsAt    time.Time     `bson:"endsAt"`
}

type WorkerClient struct {
	numWorkers     int
	productChannel chan ProductPayload
	batchSize      int
	workerWg       sync.WaitGroup

	closeOnce sync.Once

	errMu    sync.Mutex
	firstErr error
}

func InitWorkerPool() *WorkerClient {
	wc := &WorkerClient{
		numWorkers:     11,
		productChannel: make(chan ProductPayload, 200),
		batchSize:      100,
	}
	for i := 0; i < wc.numWorkers; i++ {
		wc.workerWg.Add(1)
		go wc.runWorker()
	}
	return wc
}

func (wc *WorkerClient) Send(p ProductPayload) {
	wc.productChannel <- p
}

func (wc *WorkerClient) Close() {
	wc.closeOnce.Do(func() {
		close(wc.productChannel)
	})
}

func (wc *WorkerClient) Wait() error {
	wc.workerWg.Wait()
	wc.errMu.Lock()
	defer wc.errMu.Unlock()
	return wc.firstErr
}

func (wc *WorkerClient) runWorker() {
	defer wc.workerWg.Done()
	localBuf := make([]ProductPayload, 0, wc.batchSize)
	for product := range wc.productChannel {
		localBuf = append(localBuf, product)
		if len(localBuf) >= wc.batchSize {
			wc.flush(localBuf)
			localBuf = localBuf[:0]
		}
	}
	if len(localBuf) > 0 {
		wc.flush(localBuf)
	}
}

func (wc *WorkerClient) flush(batch []ProductPayload) {
	toFlush := make([]ProductPayload, len(batch))
	copy(toFlush, batch)
	if err := wc.flushToRedis(toFlush); err != nil {
		log.Printf("worker: flush failed: %v", err)
		wc.errMu.Lock()
		if wc.firstErr == nil {
			wc.firstErr = err
		}
		wc.errMu.Unlock()
	}
}

func (wc *WorkerClient) flushToRedis(batch []ProductPayload) error {
	const maxAttempts = 3
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := wc.flushOnce(batch); err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	return fmt.Errorf("flush failed after %d attempts: %w", maxAttempts, lastErr)
}

func (wc *WorkerClient) flushOnce(batch []ProductPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := db.GetRedisClient()
	pipe := rdb.Pipeline()

	for _, p := range batch {
		priceStr := strconv.FormatFloat(p.Price, 'f', -1, 64)
		productId := p.ID.Hex()
		eventId := p.EventID.Hex()

		stockKey := rediskeys.StockKey(productId, eventId)
		pipe.HSet(ctx, stockKey, map[string]any{
			"stock": p.Frequency,
			"price": priceStr,
		})

		metaKey := rediskeys.MetaKey(productId, eventId)
		pipe.HSet(ctx, metaKey, map[string]any{
			"productId": productId,
			"eventId":   eventId,
			"sellerId":  p.SellerID.Hex(),
			"title":     p.Title,
			"currency":  p.Currency,
		})

		trackingKey := rediskeys.EventProductTrackingKey(eventId)
		pipe.SAdd(ctx, trackingKey, stockKey, metaKey)

		if !p.EndsAt.IsZero() {
			if ttl := time.Until(p.EndsAt); ttl > 0 {
				pipe.Expire(ctx, stockKey, ttl)
				pipe.Expire(ctx, metaKey, ttl)
				pipe.Expire(ctx, trackingKey, ttl)
			}
		}
	}

	_, err := pipe.Exec(ctx)
	return err
}
