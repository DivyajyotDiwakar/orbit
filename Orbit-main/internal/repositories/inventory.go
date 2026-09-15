package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"Orbit/internal/rediskeys"
	"Orbit/internal/utils"
	"Orbit/internal/worker"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InsertInventory(ctx context.Context, inventory models.Inventory) error {
	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(tCtx, inventory)
	if err != nil {
		return fmt.Errorf("insert inventory: %w", err)
	}
	return nil
}

func BulkInsertInventory(ctx context.Context, inventories []models.Inventory) error {
	if len(inventories) == 0 {
		return nil
	}

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	docs := make([]interface{}, len(inventories))
	for i, inv := range inventories {
		docs[i] = inv
	}

	opts := options.InsertMany().SetOrdered(false)
	result, err := collection.InsertMany(tCtx, docs, opts)
	if err != nil {
		var bwe mongo.BulkWriteException
		if errors.As(err, &bwe) {
			return fmt.Errorf("partial insert: %d succeeded, %d failed: %w",
				len(result.InsertedIDs), len(bwe.WriteErrors), err)
		}
		return fmt.Errorf("bulk insert: %w", err)
	}
	return nil
}

func GetAllProductsByEvent(ctx context.Context, eventIdStr string) ([]*models.Inventory, error) {
	eventId, err := utils.GetObjectFiedIdFromString(eventIdStr)
	if err != nil {
		return nil, err
	}

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(tCtx, bson.M{"eventId": eventId})
	if err != nil {
		return nil, fmt.Errorf("find products by event: %w", err)
	}
	defer cursor.Close(tCtx)

	var products []*models.Inventory
	for cursor.Next(tCtx) {
		var item models.Inventory
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("decode product: %w", err)
		}
		products = append(products, &item)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor: %w", err)
	}

	return products, nil
}

func GetProductByID(ctx context.Context, productIdStr string) (*models.Inventory, error) {
	productId, err := utils.GetObjectFiedIdFromString(productIdStr)
	if err != nil {
		return nil, err
	}

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var item models.Inventory
	err = collection.FindOne(tCtx, bson.M{"_id": productId}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	return &item, nil
}

func UpdateProduct(ctx context.Context, productIdStr string, updateData bson.M) error {
	productId, err := utils.GetObjectFiedIdFromString(productIdStr)
	if err != nil {
		return err
	}

	updateData["updatedAt"] = time.Now()

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(tCtx, bson.M{"_id": productId}, bson.M{"$set": updateData})
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func DeleteProduct(ctx context.Context, productIdStr string) error {
	productId, err := utils.GetObjectFiedIdFromString(productIdStr)
	if err != nil {
		return err
	}

	collection := db.GetInstance().Collection("inventories")
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(tCtx, bson.M{"_id": productId})
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}

func PullProducts(sellerId bson.ObjectID, eventId bson.ObjectID) error {
	pool := worker.InitWorkerPool()
	eventIdStr := eventId.Hex()

	collection := db.GetInstance().Collection("inventories")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	findOptions := options.Find().
		SetBatchSize(200).
		SetProjection(bson.M{
			"_id":       1,
			"eventId":   1,
			"sellerId":  1,
			"price":     1,
			"frequency": 1,
			"title":     1,
			"currency":  1,
			"endsAt":    1,
		})

	cursor, err := collection.Find(ctx,
		bson.M{"sellerId": sellerId, "eventId": eventId},
		findOptions,
	)
	if err != nil {
		pool.Close()
		pool.Wait()
		return fmt.Errorf("find inventory: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product worker.ProductPayload
		if err := cursor.Decode(&product); err != nil {
			log.Printf("PullProducts: decode error (skipping): %v", err)
			continue
		}
		pool.Send(product)
	}

	pool.Close()
	workerErr := pool.Wait()
	cursorErr := cursor.Err()

	if workerErr != nil || cursorErr != nil {
		if cleanupErr := CleanupEventProductKeys(context.Background(), eventIdStr); cleanupErr != nil {
			log.Printf("CRITICAL: PullProducts cleanup failed for event %s after load error: %v — "+
				"orphan product keys may remain in Redis, manual check required",
				eventIdStr, cleanupErr)
		}
	}

	if workerErr != nil {
		return fmt.Errorf("redis flush: %w", workerErr)
	}
	if cursorErr != nil {
		return fmt.Errorf("cursor: %w", cursorErr)
	}

	return nil
}

func CleanupEventProductKeys(ctx context.Context, eventId string) error {
	rdb := db.GetRedisClient()
	trackingKey := rediskeys.EventProductTrackingKey(eventId)

	keys, err := rdb.SMembers(ctx, trackingKey).Result()
	if err != nil {
		return fmt.Errorf("fetch product keys: %w", err)
	}

	if len(keys) > 0 {
		if err := rdb.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete product keys: %w", err)
		}
	}

	if err := rdb.Del(ctx, trackingKey).Err(); err != nil {
		return fmt.Errorf("delete tracking key: %w", err)
	}

	return nil
}
