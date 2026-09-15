package utils

import "go.mongodb.org/mongo-driver/v2/bson"

func GetObjectFiedIdFromString(id string) (bson.ObjectID, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return objectId, nil
}
