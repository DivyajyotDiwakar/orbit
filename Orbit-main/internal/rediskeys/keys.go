package rediskeys

import "fmt"

func StockKey(productId, eventId string) string {
	return fmt.Sprintf("product:%s:%s", productId, eventId)
}

func MetaKey(productId, eventId string) string {
	return fmt.Sprintf("productmeta:%s:%s", productId, eventId)
}

func EventProductTrackingKey(eventId string) string {
	return fmt.Sprintf("event:%s:productKeys", eventId)
}

func BookingKey(userId, productId, eventId string) string {
	return fmt.Sprintf("booking:%s:%s:%s", userId, productId, eventId)
}

func WinnersKey(eventId, productId string) string {
	return fmt.Sprintf("event:%s:product:%s:winners", eventId, productId)
}

func SaleStatusKey(eventId string) string {
	return fmt.Sprintf("event:%s:saleStatus", eventId)
}

func EventBookingsKey(eventId string) string {
	return fmt.Sprintf("event:%s:bookings", eventId)
}
