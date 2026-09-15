package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func EmailExists(ctx context.Context, email string) (bool, error) {
	collection := db.GetInstance().Collection("users")
	count, err := collection.CountDocuments(ctx, bson.M{"emailId": email})
	return count > 0, err
}

func AddAccount(ctx context.Context, user *models.User) (*models.User, error) {
	collection := db.GetInstance().Collection("users")
	_, err := collection.InsertOne(ctx, user)
	return user, err
}

func GetUserFromEmail(ctx context.Context, email string) (*models.User, error) {
	collection := db.GetInstance().Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{
		"emailId": email,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No Email Exists")
		}
		return nil, err
	}

	return &user, nil
}

func GetUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[string]models.User, error) {
	collection := db.GetInstance().Collection("users")
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make(map[string]models.User)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users[user.ID.Hex()] = user
	}

	return users, nil
}

func FindUserById(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	collection := db.GetInstance().Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user exists with that id")
		}
		return nil, err
	}
	return &user, nil
}
