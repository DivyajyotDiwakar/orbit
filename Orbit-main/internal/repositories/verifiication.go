package repositories

import (
	"Orbit/internal/db"
	"Orbit/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	VerificationStatusPending  = "pending"
	VerificationStatusApproved = "approved"
	VerificationStatusRejected = "rejected"
)

var ErrVerificationNotFound = errors.New("verification request not found")

type VerificationStateError struct {
	Current   string
	Attempted string
}

func (e *VerificationStateError) Error() string {
	return fmt.Sprintf("cannot set verification status from %s to %s", e.Current, e.Attempted)
}

func verificationCollection() *mongo.Collection {
	return db.GetInstance().Collection("verification")
}

func RegisterSellerVerificationRequest(req *models.VerificationRequest) error {
	req.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := verificationCollection().InsertOne(ctx, req)
	if err != nil {
		return fmt.Errorf("register verification request: %w", err)
	}
	return nil
}

func GetPendingVerifications(ctx context.Context) ([]models.VerificationRequest, error) {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"createdAt": 1})
	cursor, err := verificationCollection().Find(tCtx, bson.M{"status": VerificationStatusPending}, opts)
	if err != nil {
		return nil, fmt.Errorf("find pending verifications: %w", err)
	}
	defer cursor.Close(tCtx)

	var reqs []models.VerificationRequest
	if err := cursor.All(tCtx, &reqs); err != nil {
		return nil, fmt.Errorf("decode pending verifications: %w", err)
	}
	return reqs, nil
}

func GetVerificationByID(ctx context.Context, idStr string) (*models.VerificationRequest, error) {
	objId, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid verification id: %w", err)
	}

	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var v models.VerificationRequest
	err = verificationCollection().FindOne(tCtx, bson.M{"_id": objId}).Decode(&v)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrVerificationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find verification: %w", err)
	}
	return &v, nil
}

func ApproveVerification(ctx context.Context, verificationIdStr string, adminId bson.ObjectID) error {
	objId, err := bson.ObjectIDFromHex(verificationIdStr)
	if err != nil {
		return fmt.Errorf("invalid verification id: %w", err)
	}

	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var v models.VerificationRequest
	if err := verificationCollection().FindOne(tCtx, bson.M{"_id": objId}).Decode(&v); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrVerificationNotFound
		}
		return fmt.Errorf("find verification: %w", err)
	}

	if v.Status == VerificationStatusApproved {
		return syncSellerVerificationDetails(tCtx, &v)
	}
	if v.Status != VerificationStatusPending {
		return &VerificationStateError{Current: v.Status, Attempted: VerificationStatusApproved}
	}

	now := time.Now()
	res, err := verificationCollection().UpdateOne(tCtx,
		bson.M{"_id": objId, "status": VerificationStatusPending},
		bson.M{"$set": bson.M{
			"status":     VerificationStatusApproved,
			"reviewedBy": adminId,
			"reviewedAt": now,
		}},
	)
	if err != nil {
		return fmt.Errorf("update verification status: %w", err)
	}
	if res.MatchedCount == 0 {
		return &VerificationStateError{Current: "unknown (concurrent update)", Attempted: VerificationStatusApproved}
	}

	return syncSellerVerificationDetails(tCtx, &v)
}

func syncSellerVerificationDetails(ctx context.Context, verification *models.VerificationRequest) error {
	collection := db.GetInstance().Collection("users")
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": verification.SellerID},
		bson.M{"$set": bson.M{
			"sellerInfo.isApproved":                    true,
			"sellerInfo.companyName":                   verification.CompanyName,
			"sellerInfo.gstin":                         verification.GSTIN,
			"sellerInfo.bankDetails.accountHolderName": verification.AccountHolderName,
			"sellerInfo.bankDetails.accountNumber":     verification.AccountNumber,
			"sellerInfo.bankDetails.ifscCode":          verification.IFSCCode,
			"sellerInfo.bankDetails.bankName":          verification.BankName,
			"updatedAt":                                time.Now(),
		}},
	)
	if err != nil {
		log.Printf("ApproveVerification: user details sync failed for seller %s: %v — "+
			"verification.status is already APPROVED (source of truth); retry ApproveVerification to re-sync",
			verification.SellerID.Hex(), err)
		return fmt.Errorf("sync seller verification details: %w", err)
	}
	return nil
}

func RejectVerification(ctx context.Context, verificationIdStr string, adminId bson.ObjectID, reason string) error {
	objId, err := bson.ObjectIDFromHex(verificationIdStr)
	if err != nil {
		return fmt.Errorf("invalid verification id: %w", err)
	}

	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var v models.VerificationRequest
	if err := verificationCollection().FindOne(tCtx, bson.M{"_id": objId}).Decode(&v); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrVerificationNotFound
		}
		return fmt.Errorf("find verification: %w", err)
	}

	if v.Status == VerificationStatusRejected {
		return nil
	}
	if v.Status != VerificationStatusPending {
		return &VerificationStateError{Current: v.Status, Attempted: VerificationStatusRejected}
	}

	now := time.Now()
	res, err := verificationCollection().UpdateOne(tCtx,
		bson.M{"_id": objId, "status": VerificationStatusPending},
		bson.M{"$set": bson.M{
			"status":          VerificationStatusRejected,
			"reviewedBy":      adminId,
			"reviewedAt":      now,
			"rejectionReason": reason,
		}},
	)
	if err != nil {
		return fmt.Errorf("update verification status: %w", err)
	}
	if res.MatchedCount == 0 {
		return &VerificationStateError{Current: "unknown (concurrent update)", Attempted: VerificationStatusRejected}
	}

	return nil
}
