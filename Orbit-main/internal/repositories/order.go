package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateOrder(ctx context.Context, order models.Order) error {
	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(tCtx, order)
	if err != nil {
		var we mongo.WriteException
		if errors.As(err, &we) {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 { // duplicate key
					return nil
				}
			}
		}
		return fmt.Errorf("create order: %w", err)
	}
	return nil
}

func BulkUpsertOrders(ctx context.Context, orders []models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	models := make([]mongo.WriteModel, len(orders))
	for i, o := range orders {
		models[i] = mongo.NewUpdateOneModel().
			SetFilter(bson.M{"reservationId": o.ReservationID}).
			SetUpdate(bson.M{"$setOnInsert": o}).
			SetUpsert(true)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := collection.BulkWrite(tCtx, models, opts)
	if err != nil {
		return fmt.Errorf("bulk upsert orders: %w", err)
	}
	return nil
}

func GetOrdersByUser(ctx context.Context, userId bson.ObjectID) ([]models.Order, error) {
	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := collection.Find(tCtx, bson.M{"userId": userId}, opts)
	if err != nil {
		return nil, fmt.Errorf("get orders by user: %w", err)
	}
	defer cursor.Close(tCtx)

	var orders []models.Order
	if err := cursor.All(tCtx, &orders); err != nil {
		return nil, fmt.Errorf("decode orders: %w", err)
	}
	return orders, nil
}

func GetOrdersByEvent(ctx context.Context, eventId bson.ObjectID) ([]models.Order, error) {
	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := collection.Find(tCtx, bson.M{"eventId": eventId}, opts)
	if err != nil {
		return nil, fmt.Errorf("get orders by event: %w", err)
	}
	defer cursor.Close(tCtx)

	var orders []models.Order
	if err := cursor.All(tCtx, &orders); err != nil {
		return nil, fmt.Errorf("decode event orders: %w", err)
	}
	return orders, nil
}

type EventAnalytics struct {
	TotalOrders     int            `json:"totalOrders"`
	TotalRevenue    float64        `json:"totalRevenue"`
	StatusBreakdown map[string]int `json:"statusBreakdown"`
}

func GetEventAnalytics(ctx context.Context, eventId bson.ObjectID) (*EventAnalytics, error) {
	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"eventId": eventId}}},
		{{Key: "$group", Value: bson.M{
			"_id":          "$status",
			"count":        bson.M{"$sum": 1},
			"totalRevenue": bson.M{"$sum": "$price"},
		}}},
	}

	cursor, err := collection.Aggregate(tCtx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("analytics aggregate: %w", err)
	}
	defer cursor.Close(tCtx)

	analytics := &EventAnalytics{
		StatusBreakdown: make(map[string]int),
	}

	var rows []struct {
		Status  string  `bson:"_id"`
		Count   int     `bson:"count"`
		Revenue float64 `bson:"totalRevenue"`
	}
	if err := cursor.All(tCtx, &rows); err != nil {
		return nil, fmt.Errorf("decode analytics: %w", err)
	}

	for _, row := range rows {
		analytics.TotalOrders += row.Count
		analytics.TotalRevenue += row.Revenue
		analytics.StatusBreakdown[row.Status] = row.Count
	}

	return analytics, nil
}
