package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderStatus string

const (
	OrderStatusConfirmed      OrderStatus = "CONFIRMED"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	ID            bson.ObjectID `bson:"_id"           json:"id"`
	UserID        bson.ObjectID `bson:"userId"         json:"userId"`
	ProductID     bson.ObjectID `bson:"productId"      json:"productId"`
	EventID       bson.ObjectID `bson:"eventId"        json:"eventId"`
	ReservationID string        `bson:"reservationId"  json:"reservationId"`
	Price         float64       `bson:"price"          json:"price"`
	Status        OrderStatus   `bson:"status"         json:"status"`
	CreatedAt     time.Time     `bson:"createdAt"      json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt"      json:"updatedAt"`
}
