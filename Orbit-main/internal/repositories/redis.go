package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/rediskeys"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type BookingResult int

const (
	BookingSoldOut        BookingResult = 0
	BookingSuccess        BookingResult = 1
	BookingAlreadyDone    BookingResult = 2
	BookingProductMissing BookingResult = 3
	BookingSaleNotActive  BookingResult = 4
)

var ErrProductNotFound = errors.New("product not found in inventory")

type Reservation struct {
	ReservationID string
	ProductID     string
	EventID       string
	UserID        string
	Price         float64
	Status        string
	BookedAt      time.Time
	ExpiresAt     time.Time
}

var reserveScript = redis.NewScript(`
	local saleStatus = redis.call('GET', KEYS[4])
	if saleStatus ~= 'LIVE' then
		return 4
	end

	if redis.call('EXISTS', KEYS[2]) == 1 then
		return 2
	end

	local fields = redis.call('HMGET', KEYS[1], 'stock', 'price')
	local stock = tonumber(fields[1])
	local price = fields[2]

	if not stock and not price then
		return 3
	end

	if not stock or stock <= 0 then
		return 0
	end

	redis.call('HINCRBY', KEYS[1], 'stock', -1)

	redis.call('HSET', KEYS[2],
		'userId',    ARGV[1],
		'productId', ARGV[2],
		'eventId',   ARGV[3],
		'price',     price,
		'status',    ARGV[6],
		'bookedAt',  ARGV[4]
	)

	local ttl = tonumber(ARGV[5])
	if ttl > 0 then
		redis.call('EXPIRE', KEYS[2], ttl)
	end

	redis.call('SADD', KEYS[3], ARGV[1])
	redis.call('SADD', KEYS[5], KEYS[2])

	-- Returning price here keeps a successful purchase to one Redis round trip.
	-- A separate HGET may time out under a burst after stock was decremented.
	return {1, price}
`)

func ReserveProduct(
	ctx context.Context,
	productId string,
	eventId string,
	userId string,
	status string,
	ttl time.Duration,
) (*Reservation, BookingResult, error) {
	rdb := db.GetRedisClient()

	keys := []string{
		rediskeys.StockKey(productId, eventId),
		rediskeys.BookingKey(userId, productId, eventId),
		rediskeys.WinnersKey(eventId, productId),
		rediskeys.SaleStatusKey(eventId),
		rediskeys.EventBookingsKey(eventId),
	}

	now := time.Now()
	args := []interface{}{
		userId,
		productId,
		eventId,
		fmt.Sprintf("%d", now.Unix()),
		fmt.Sprintf("%d", int64(ttl.Seconds())),
		status,
	}

	raw, err := reserveScript.Run(ctx, rdb, keys, args...).Result()
	if err != nil {
		return nil, BookingSoldOut, fmt.Errorf("reserve script: %w", err)
	}

	result, priceStr, err := parseReservationResult(raw)
	if err != nil {
		return nil, BookingSoldOut, err
	}
	if result != BookingSuccess {
		return nil, result, nil
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, BookingSuccess, fmt.Errorf("parse price %q: %w", priceStr, err)
	}

	reservation := &Reservation{
		ReservationID: keys[1],
		ProductID:     productId,
		EventID:       eventId,
		UserID:        userId,
		Price:         price,
		Status:        status,
		BookedAt:      now,
	}
	if ttl > 0 {
		reservation.ExpiresAt = now.Add(ttl)
	}

	return reservation, BookingSuccess, nil
}

func parseReservationResult(raw interface{}) (BookingResult, string, error) {
	switch value := raw.(type) {
	case int64:
		return BookingResult(value), "", nil
	case []interface{}:
		if len(value) != 2 {
			return BookingSoldOut, "", fmt.Errorf("reserve script: unexpected result shape %v", raw)
		}
		code, ok := value[0].(int64)
		if !ok {
			return BookingSoldOut, "", fmt.Errorf("reserve script: invalid result code %T", value[0])
		}
		price, ok := value[1].(string)
		if !ok {
			return BookingSoldOut, "", fmt.Errorf("reserve script: invalid price %T", value[1])
		}
		return BookingResult(code), price, nil
	default:
		return BookingSoldOut, "", fmt.Errorf("reserve script: unexpected result type %T", raw)
	}
}
