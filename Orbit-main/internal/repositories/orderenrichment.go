package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetProductTitlesByIDs(ctx context.Context, ids []bson.ObjectID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1})
	cursor, err := collection.Find(tCtx, bson.M{"_id": bson.M{"$in": ids}}, opts)
	if err != nil {
		return nil, fmt.Errorf("get product titles: %w", err)
	}
	defer cursor.Close(tCtx)

	var rows []struct {
		ID    bson.ObjectID `bson:"_id"`
		Title string        `bson:"title"`
	}
	if err := cursor.All(tCtx, &rows); err != nil {
		return nil, fmt.Errorf("decode product titles: %w", err)
	}

	for _, r := range rows {
		result[r.ID.Hex()] = r.Title
	}
	return result, nil
}

func GetEventNamesByIDs(ctx context.Context, ids []bson.ObjectID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	collection := db.GetInstance().Collection("events")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{"_id": 1, "eventName": 1})
	cursor, err := collection.Find(tCtx, bson.M{"_id": bson.M{"$in": ids}}, opts)
	if err != nil {
		return nil, fmt.Errorf("get event names: %w", err)
	}
	defer cursor.Close(tCtx)

	var rows []struct {
		ID        bson.ObjectID `bson:"_id"`
		EventName string        `bson:"eventName"`
	}
	if err := cursor.All(tCtx, &rows); err != nil {
		return nil, fmt.Errorf("decode event names: %w", err)
	}

	for _, r := range rows {
		result[r.ID.Hex()] = r.EventName
	}
	return result, nil
}

func GetSellerEventIDs(ctx context.Context, sellerId bson.ObjectID) ([]bson.ObjectID, error) {
	collection := db.GetInstance().Collection("events")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := collection.Find(tCtx, bson.M{"sellerId": sellerId}, opts)
	if err != nil {
		return nil, fmt.Errorf("get seller event ids: %w", err)
	}
	defer cursor.Close(tCtx)

	var rows []struct {
		ID bson.ObjectID `bson:"_id"`
	}
	if err := cursor.All(tCtx, &rows); err != nil {
		return nil, fmt.Errorf("decode seller event ids: %w", err)
	}

	ids := make([]bson.ObjectID, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}
	return ids, nil
}

func GetOrdersBySeller(ctx context.Context, sellerId bson.ObjectID) ([]models.Order, error) {
	eventIds, err := GetSellerEventIDs(ctx, sellerId)
	if err != nil {
		return nil, err
	}
	if len(eventIds) == 0 {
		return []models.Order{}, nil
	}

	collection := db.GetInstance().Collection("orders")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := collection.Find(tCtx, bson.M{"eventId": bson.M{"$in": eventIds}}, opts)
	if err != nil {
		return nil, fmt.Errorf("get orders by seller: %w", err)
	}
	defer cursor.Close(tCtx)

	var orders []models.Order
	if err := cursor.All(tCtx, &orders); err != nil {
		return nil, fmt.Errorf("decode seller orders: %w", err)
	}
	return orders, nil
}
