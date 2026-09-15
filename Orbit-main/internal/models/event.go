package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Event struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SellerID    bson.ObjectID `bson:"sellerId" json:"sellerId"`
	EventName   string        `bson:"eventName" json:"eventName"`
	Description string        `bson:"description" json:"description"`
	ScheduledAt time.Time     `bson:"scheduledAt" json:"scheduledAt"`
	IsLive      bool          `bson:"isLive" json:"isLive"`
	ImageBanner string        `bson:"imageBanner" json:"imageBanner"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
