package repositories

import (
	"Orbit/internal/db"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LiveEvent struct {
	ID          bson.ObjectID `bson:"_id"`
	EventName   string        `bson:"eventName"`
	Description string        `bson:"description"`
	ImageBanner string        `bson:"imageBanner"`
	ScheduledAt time.Time     `bson:"scheduledAt"`
}

type BuyerProduct struct {
	ID          bson.ObjectID `bson:"_id"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Price       float64       `bson:"price"`
	Currency    string        `bson:"currency"`
	Image       string        `bson:"image"`
}

type BuyerProductWithStock struct {
	BuyerProduct
	AvailableStock int
}

func GetLiveEvents(ctx context.Context) ([]LiveEvent, error) {
	collection := db.GetInstance().Collection("events")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{
		"_id":         1,
		"eventName":   1,
		"description": 1,
		"imageBanner": 1,
		"scheduledAt": 1,
	})

	cursor, err := collection.Find(tCtx, bson.M{"isLive": true}, opts)
	if err != nil {
		return nil, fmt.Errorf("get live events: %w", err)
	}
	defer cursor.Close(tCtx)

	var events []LiveEvent
	if err := cursor.All(tCtx, &events); err != nil {
		return nil, fmt.Errorf("decode live events: %w", err)
	}

	return events, nil
}

func GetEventProductsWithStock(ctx context.Context, eventId bson.ObjectID) ([]BuyerProductWithStock, error) {
	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{
		"_id":         1,
		"title":       1,
		"description": 1,
		"price":       1,
		"currency":    1,
		"image":       1,
	})

	cursor, err := collection.Find(tCtx, bson.M{"eventId": eventId}, opts)
	if err != nil {
		return nil, fmt.Errorf("find event products: %w", err)
	}
	defer cursor.Close(tCtx)

	var products []BuyerProduct
	if err := cursor.All(tCtx, &products); err != nil {
		return nil, fmt.Errorf("decode event products: %w", err)
	}

	if len(products) == 0 {
		return []BuyerProductWithStock{}, nil
	}

	rdb := db.GetRedisClient()
	stockCtx, stockCancel := context.WithTimeout(ctx, 2*time.Second)
	defer stockCancel()

	pipe := rdb.Pipeline()
	eventIdStr := eventId.Hex()

	stockCmds := make([]*redis.StringCmd, len(products))
	for i, p := range products {
		stockKey := fmt.Sprintf("product:%s:%s", p.ID.Hex(), eventIdStr)
		stockCmds[i] = pipe.HGet(stockCtx, stockKey, "stock")
	}

	pipe.Exec(stockCtx)

	result := make([]BuyerProductWithStock, len(products))
	for i, p := range products {
		stock := 0
		if val, err := stockCmds[i].Result(); err == nil {
			stock, _ = strconv.Atoi(val)
		}

		result[i] = BuyerProductWithStock{
			BuyerProduct:   p,
			AvailableStock: stock,
		}
	}

	return result, nil
}
