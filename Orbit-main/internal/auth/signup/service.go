package signup

import (
	"Orbit/internal/models"
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserService(ctx context.Context, req RegisterRequest) (string, *models.User, error) {
	exists, err := repositories.EmailExists(ctx, req.EmailId)
	if err != nil {
		return "", nil, err
	}
	if exists {
		return "", nil, errors.New("Email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return "", nil, errors.New("Password hashing failed")
	}

	now := time.Now()
	user := &models.User{
		ID:              bson.NewObjectID(),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		EmailId:         req.EmailId,
		Age:             req.Age,
		Password:        string(hash),
		Role:            req.Role,
		IsEmailVerified: true, // should be false and verified in production environment.
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if req.Role == "buyer" {
		user.BuyerInfo = &models.BuyerProfile{
			IsApproved: true,
		}
	} else if req.Role == "seller" {
		user.SellerInfo = &models.SellerProfile{
			IsApproved: true, // will keep false in production environment
		}
	}

	_, err = repositories.AddAccount(ctx, user)
	if err != nil {
		return "", nil, errors.New("Database insert failed")
	}

	token, err := utils.CreateJwtToken(*user)
	if err != nil {
		return "", user, errors.New("token creation failed")
	}

	return token, user, nil
}
