package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VerificationRequest struct {
	ID                bson.ObjectID `bson:"_id,omitempty"          json:"id"`
	SellerID          bson.ObjectID `bson:"sellerId"                json:"sellerId"`
	EmailId           string        `bson:"emailId"                 json:"emailId"`
	CompanyName       string        `bson:"companyName"             json:"companyName"`
	GSTIN             string        `bson:"gstin"                   json:"gstin"`
	AccountHolderName string        `bson:"accountHolderName"       json:"accountHolderName"`
	AccountNumber     string        `bson:"accountNumber"           json:"accountNumber"`
	IFSCCode          string        `bson:"ifscCode"                json:"ifscCode"`
	BankName          string        `bson:"bankName"                json:"bankName"`
	Status            string        `bson:"status"                  json:"status"` // "pending" | "approved" | "rejected"

	CreatedAt       time.Time      `bson:"createdAt"                 json:"createdAt"`
	ReviewedBy      *bson.ObjectID `bson:"reviewedBy,omitempty"      json:"reviewedBy,omitempty"`
	ReviewedAt      *time.Time     `bson:"reviewedAt,omitempty"      json:"reviewedAt,omitempty"`
	RejectionReason string         `bson:"rejectionReason,omitempty" json:"rejectionReason,omitempty"`
}
