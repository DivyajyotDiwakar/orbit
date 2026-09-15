package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"Orbit/internal/rediskeys"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SaleStatus string

const (
	SaleStatusNotStarted SaleStatus = "NOT_STARTED"
	SaleStatusLive       SaleStatus = "LIVE"
	SaleStatusPaused     SaleStatus = "PAUSED"
	SaleStatusStopped    SaleStatus = "STOPPED"
)

type SaleTransitionError struct {
	Current   SaleStatus
	Attempted SaleStatus
}

func (e *SaleTransitionError) Error() string {
	return fmt.Sprintf("cannot transition sale from %s to %s", e.Current, e.Attempted)
}

var transitionScript = redis.NewScript(`
	local current = redis.call('GET', KEYS[1])
	if current == false then
		current = 'NOT_STARTED'
	end

	local target = ARGV[#ARGV]

	for i = 1, #ARGV - 1 do
		if ARGV[i] == current then
			redis.call('SET', KEYS[1], target)
			return {1, current}
		end
	end

	return {0, current}
`)

func transitionSaleStatus(
	ctx context.Context,
	eventId string,
	allowedFrom []SaleStatus,
	target SaleStatus,
) (bool, SaleStatus, error) {
	rdb := db.GetRedisClient()

	args := make([]interface{}, 0, len(allowedFrom)+1)
	for _, s := range allowedFrom {
		args = append(args, string(s))
	}
	args = append(args, string(target))

	raw, err := transitionScript.Run(ctx, rdb, []string{rediskeys.SaleStatusKey(eventId)}, args...).Result()
	if err != nil {
		return false, "", fmt.Errorf("transition script: %w", err)
	}

	result, ok := raw.([]interface{})
	if !ok || len(result) != 2 {
		return false, "", fmt.Errorf("transition script: unexpected result shape %v", raw)
	}

	success, _ := result[0].(int64)
	current, _ := result[1].(string)

	return success == 1, SaleStatus(current), nil
}

func GetSaleStatus(ctx context.Context, eventId string) (SaleStatus, error) {
	val, err := db.GetRedisClient().Get(ctx, rediskeys.SaleStatusKey(eventId)).Result()
	if errors.Is(err, redis.Nil) {
		return SaleStatusNotStarted, nil
	}
	if err != nil {
		return "", fmt.Errorf("get sale status: %w", err)
	}
	return SaleStatus(val), nil
}

func LiveSale(ctx context.Context, eventId string, sellerId bson.ObjectID, eventObjId bson.ObjectID) error {
	ok, current, err := transitionSaleStatus(ctx, eventId, []SaleStatus{SaleStatusNotStarted}, SaleStatusLive)
	if err != nil {
		return fmt.Errorf("live sale: %w", err)
	}
	if !ok {
		return &SaleTransitionError{Current: current, Attempted: SaleStatusLive}
	}

	if err := PullProducts(sellerId, eventObjId); err != nil {
		rollbackLiveTransition(ctx, eventId, "PullProducts failed")
		return fmt.Errorf("load inventory: %w", err)
	}

	if err := updateEventLiveStatus(ctx, eventObjId, true); err != nil {
		rollbackLiveTransition(ctx, eventId, "mongo isLive update failed")
		return fmt.Errorf("update event live status: %w", err)
	}

	return nil
}

func rollbackLiveTransition(ctx context.Context, eventId string, reason string) {
	ok, current, err := transitionSaleStatus(ctx, eventId, []SaleStatus{SaleStatusLive}, SaleStatusNotStarted)
	if err != nil {
		log.Printf("CRITICAL: LiveSale rollback failed for event %s (%s): %v — "+
			"status may be stuck LIVE with incomplete setup, manual check required",
			eventId, reason, err)
		return
	}
	if !ok {
		log.Printf("LiveSale rollback skipped for event %s (%s): status already moved to %s concurrently",
			eventId, reason, current)
	}
}

func PauseSale(ctx context.Context, eventId string) error {
	ok, current, err := transitionSaleStatus(ctx, eventId, []SaleStatus{SaleStatusLive}, SaleStatusPaused)
	if err != nil {
		return fmt.Errorf("pause sale: %w", err)
	}
	if !ok {
		return &SaleTransitionError{Current: current, Attempted: SaleStatusPaused}
	}
	return nil
}

func ResumeSale(ctx context.Context, eventId string) error {
	ok, current, err := transitionSaleStatus(ctx, eventId, []SaleStatus{SaleStatusPaused}, SaleStatusLive)
	if err != nil {
		return fmt.Errorf("resume sale: %w", err)
	}
	if !ok {
		return &SaleTransitionError{Current: current, Attempted: SaleStatusLive}
	}
	return nil
}

func StopSale(ctx context.Context, eventId string, eventObjId bson.ObjectID) error {
	ok, current, err := transitionSaleStatus(
		ctx, eventId,
		[]SaleStatus{SaleStatusLive, SaleStatusPaused},
		SaleStatusStopped,
	)
	if err != nil {
		return fmt.Errorf("stop sale: %w", err)
	}

	if !ok {
		if current == SaleStatusStopped {
			return syncStoppedSaleOrders(ctx, eventId, eventObjId)
		}
		return &SaleTransitionError{Current: current, Attempted: SaleStatusStopped}
	}

	return syncStoppedSaleOrders(ctx, eventId, eventObjId)
}

func syncStoppedSaleOrders(ctx context.Context, eventId string, eventObjId bson.ObjectID) error {
	rdb := db.GetRedisClient()

	bookingKeys, err := rdb.SMembers(ctx, rediskeys.EventBookingsKey(eventId)).Result()
	if err != nil {
		return fmt.Errorf("fetch booking keys: %w", err)
	}

	if len(bookingKeys) > 0 {
		pipe := rdb.Pipeline()
		cmds := make([]*redis.MapStringStringCmd, len(bookingKeys))
		for i, key := range bookingKeys {
			cmds[i] = pipe.HGetAll(ctx, key)
		}

		if _, execErr := pipe.Exec(ctx); execErr != nil {
			log.Printf("StopSale: pipeline HGETALL error for event %s (checking individual keys next): %v",
				eventId, execErr)
		}

		orders := make([]models.Order, 0, len(bookingKeys))
		var skipped int
		for i, cmd := range cmds {
			data, err := cmd.Result()
			if err != nil || len(data) == 0 {
				log.Printf("StopSale: skip booking key %s: %v", bookingKeys[i], err)
				skipped++
				continue
			}
			order, err := bookingHashToOrder(data, bookingKeys[i])
			if err != nil {
				log.Printf("StopSale: parse booking %s: %v", bookingKeys[i], err)
				skipped++
				continue
			}
			orders = append(orders, order)
		}

		if skipped > 0 {
			log.Printf("StopSale: event %s — %d of %d bookings skipped during sync",
				eventId, skipped, len(bookingKeys))
		}

		if err := BulkUpsertOrders(ctx, orders); err != nil {
			return fmt.Errorf("sync orders to mongo: %w", err)
		}
	}

	if err := updateEventLiveStatus(ctx, eventObjId, false); err != nil {
		log.Printf("StopSale: update event isLive failed (non-fatal) for event %s: %v", eventId, err)
	}

	go cleanupEventRedisKeys(eventId, bookingKeys)

	return nil
}

func bookingHashToOrder(data map[string]string, reservationId string) (models.Order, error) {
	userId, err := bson.ObjectIDFromHex(data["userId"])
	if err != nil {
		return models.Order{}, fmt.Errorf("parse userId: %w", err)
	}
	productId, err := bson.ObjectIDFromHex(data["productId"])
	if err != nil {
		return models.Order{}, fmt.Errorf("parse productId: %w", err)
	}
	eventId, err := bson.ObjectIDFromHex(data["eventId"])
	if err != nil {
		return models.Order{}, fmt.Errorf("parse eventId: %w", err)
	}
	price, err := strconv.ParseFloat(data["price"], 64)
	if err != nil {
		return models.Order{}, fmt.Errorf("parse price: %w", err)
	}

	bookedAtUnix, _ := strconv.ParseInt(data["bookedAt"], 10, 64)
	bookedAt := time.Unix(bookedAtUnix, 0)

	now := time.Now()
	return models.Order{
		ID:            bson.NewObjectID(),
		UserID:        userId,
		ProductID:     productId,
		EventID:       eventId,
		ReservationID: reservationId,
		Price:         price,
		Status:        models.OrderStatus(data["status"]),
		CreatedAt:     bookedAt,
		UpdatedAt:     now,
	}, nil
}

func updateEventLiveStatus(ctx context.Context, eventId bson.ObjectID, isLive bool) error {
	col := db.GetInstance().Collection("events")
	_, err := col.UpdateOne(ctx,
		bson.M{"_id": eventId},
		bson.M{"$set": bson.M{"isLive": isLive, "updatedAt": time.Now()}},
	)
	return err
}

func cleanupEventRedisKeys(eventId string, bookingKeys []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rdb := db.GetRedisClient()

	if len(bookingKeys) > 0 {
		if err := rdb.Del(ctx, bookingKeys...).Err(); err != nil {
			log.Printf("cleanup: del booking keys: %v", err)
		}
	}

	if err := rdb.Del(ctx, rediskeys.EventBookingsKey(eventId)).Err(); err != nil {
		log.Printf("cleanup: del event bookings key: %v", err)
	}

	if err := CleanupEventProductKeys(ctx, eventId); err != nil {
		log.Printf("cleanup: del product keys: %v", err)
	}

	if err := rdb.Expire(ctx, rediskeys.SaleStatusKey(eventId), 30*24*time.Hour).Err(); err != nil {
		log.Printf("cleanup: expire sale status key: %v", err)
	}
}
