package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Inventory struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SellerID bson.ObjectID `bson:"sellerId" json:"sellerId"`
	EventID  bson.ObjectID `bson:"eventId" json:"eventId"`

	Title       string  `bson:"title" json:"title"`
	Description string  `bson:"description" json:"description"`
	Price       float64 `bson:"price" json:"price"`
	Frequency   int     `bson:"frequency" json:"frequency"`
	Image       string  `bson:"image" json:"image"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
