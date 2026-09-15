package sale

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func LiveSaleService(sellerId string, eventId string) error {
	sellerObjId, err := utils.GetObjectFiedIdFromString(sellerId)
	if err != nil {
		return err
	}
	eventObjId, err := utils.GetObjectFiedIdFromString(eventId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return repositories.LiveSale(ctx, eventId, sellerObjId, eventObjId)
}

func PauseSaleService(sellerId string, eventId string) error {
	if _, err := utils.GetObjectFiedIdFromString(sellerId); err != nil {
		return err
	}
	ctx := context.Background()
	return repositories.PauseSale(ctx, eventId)
}

func ResumeSaleService(sellerId string, eventId string) error {
	if _, err := utils.GetObjectFiedIdFromString(sellerId); err != nil {
		return err
	}
	ctx := context.Background()
	return repositories.ResumeSale(ctx, eventId)
}

func StopSaleService(sellerId string, eventId string) error {
	if _, err := utils.GetObjectFiedIdFromString(sellerId); err != nil {
		return err
	}
	eventObjId, err := bson.ObjectIDFromHex(eventId)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return repositories.StopSale(ctx, eventId, eventObjId)
}
