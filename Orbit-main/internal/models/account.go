package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type SellerProfile struct {
	CompanyName     string    `bson:"companyName,omitempty" json:"companyName,omitempty"`
	GSTIN           string    `bson:"gstin,omitempty" json:"gstin,omitempty"`
	IsApproved      bool      `bson:"isApproved" json:"isApproved"`
	BankDetails     BankInfo  `bson:"bankDetails,omitempty" json:"-"`
	LastRequestedAt time.Time `bson:"lastRequestedAt,omitempty" json:"lastRequestedAt,omitempty"`
	NextAllowedAt   time.Time `bson:"nextAllowedAt,omitempty" json:"nextAllowedAt,omitempty"`
}
type BankInfo struct {
	AccountHolderName string `bson:"accountHolderName,omitempty" json:"accountHolderName,omitempty"`
	AccountNumber     string `bson:"accountNumber,omitempty" json:"accountNumber,omitempty"`
	IFSCCode          string `bson:"ifscCode,omitempty" json:"ifscCode,omitempty"`
	BankName          string `bson:"bankName,omitempty" json:"bankName,omitempty"`
}

type BuyerProfile struct {
	PhoneNumber string   `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
	Addresses   []string `bson:"addresses,omitempty" json:"addresses,omitempty"`
	IsApproved  bool     `bson:"isApproved" json:"isApproved"`
}

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	FirstName string        `bson:"firstName" json:"firstName"`
	LastName  string        `bson:"lastName,omitempty" json:"lastName,omitempty"`
	EmailId   string        `bson:"emailId" json:"emailId"`
	Password  string        `bson:"password" json:"-"`
	Age       int           `bson:"age,omitempty" json:"age,omitempty"`
	Role      string        `bson:"role" json:"role"` // enum: "admin", "buyer", "seller", "user"
	Image     string        `bson:"image,omitempty" json:"image,omitempty"`

	IsEmailVerified bool `bson:"isEmailVerified" json:"isEmailVerified"`
	IsActive        bool `bson:"isActive" json:"isActive"` // to block if needed.

	BuyerInfo  *BuyerProfile  `bson:"buyerInfo,omitempty" json:"buyerInfo,omitempty"`
	SellerInfo *SellerProfile `bson:"sellerInfo,omitempty" json:"sellerInfo,omitempty"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
