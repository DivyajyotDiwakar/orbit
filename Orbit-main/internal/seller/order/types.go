package order

import (
	"time"
)

type SellerOrderView struct {
	OrderID       string    `json:"orderId"`
	ProductID     string    `json:"productId"`
	ProductTitle  string    `json:"productTitle"`
	EventID       string    `json:"eventId"`
	EventName     string    `json:"eventName"`
	BuyerID       string    `json:"buyerId"`
	ReservationID string    `json:"reservationId"`
	CustomerName  string    `json:"customerName"`
	CustomerEmail string    `json:"customerEmail"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	BookedAt      time.Time `json:"bookedAt"`
}
