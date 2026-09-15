package product

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func RegisterProductService(
	ctx context.Context,
	sellerIdStr string,
	eventIdStr string,
	req RegisterProductRequestBody,
	imageURL string,
) error {
	sellerId, err := utils.GetObjectFiedIdFromString(sellerIdStr)
	if err != nil {
		return ErrInvalidSeller
	}

	eventId, err := utils.GetObjectFiedIdFromString(eventIdStr)
	if err != nil {
		return ErrInvalidEvent
	}

	event, err := repositories.GetAnEvent(ctx, eventIdStr)
	if err != nil || event == nil {
		return ErrEventNotFound
	}

	if event.SellerID != sellerId {
		return ErrUnauthorized
	}

	now := time.Now()
	inventory := models.Inventory{
		ID:          bson.NewObjectID(),
		SellerID:    sellerId,
		EventID:     eventId,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Frequency:   req.Frequency,
		Image:       imageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return repositories.InsertInventory(ctx, inventory)
}

func GetAllEventProductsService(ctx context.Context, eventIdStr string) ([]*models.Inventory, error) {
	return repositories.GetAllProductsByEvent(ctx, eventIdStr)
}

func GetAnEventProductService(ctx context.Context, productIdStr string) (*models.Inventory, error) {
	product, err := repositories.GetProductByID(ctx, productIdStr)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func UpdateAnEventProductService(
	ctx context.Context,
	sellerIdStr string,
	eventIdStr string,
	productIdStr string,
	req UpdateProductRequestBody,
	newImageURL string,
) error {
	sellerId, err := utils.GetObjectFiedIdFromString(sellerIdStr)
	if err != nil {
		return ErrInvalidSeller
	}

	event, err := repositories.GetAnEvent(ctx, eventIdStr)
	if err != nil || event == nil {
		return ErrEventNotFound
	}

	if event.SellerID != sellerId {
		return ErrUnauthorized
	}

	product, err := repositories.GetProductByID(ctx, productIdStr)
	if err != nil {
		return err // real DB error
	}
	if product == nil {
		return ErrProductNotFound
	}

	updateData := bson.M{}
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.Price != nil {
		updateData["price"] = *req.Price
	}
	if req.Frequency != nil {
		updateData["frequency"] = *req.Frequency
	}
	if newImageURL != "" {
		updateData["image"] = newImageURL
	}

	if len(updateData) == 0 {
		return ErrNoUpdateFields
	}

	return repositories.UpdateProduct(ctx, productIdStr, updateData)
}

func DeleteAnEventProductService(
	ctx context.Context,
	sellerIdStr string,
	eventIdStr string,
	productIdStr string,
) error {
	sellerId, err := utils.GetObjectFiedIdFromString(sellerIdStr)
	if err != nil {
		return ErrInvalidSeller
	}

	event, err := repositories.GetAnEvent(ctx, eventIdStr)
	if err != nil || event == nil {
		return ErrEventNotFound
	}

	if event.SellerID != sellerId {
		return ErrUnauthorized
	}

	product, err := repositories.GetProductByID(ctx, productIdStr)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	return repositories.DeleteProduct(ctx, productIdStr)
}
