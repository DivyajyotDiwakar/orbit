package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	SellerID       string
	BuyerID        string
	SellerEmail    string
	BuyerEmail     string
	SellerPassword = "Password123"
	BuyerPassword  = "Password123"
)

type AuthUser struct {
	ID        string `json:"_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	EmailId   string `json:"emailId"`
	Age       int    `json:"age,omitempty"`
	Role      string `json:"role"`
}

type AuthResponse struct {
	Message string   `json:"message"`
	User    AuthUser `json:"user"`
}

func AuthFlow() {
	fmt.Println("========== AUTH FLOW ==========")

	BuyerEmail = fmt.Sprintf("buyer-%d@example.com", time.Now().UnixNano())
	SellerEmail = fmt.Sprintf("seller-%d@example.com", time.Now().UnixNano())

	BuyerSignup()
	BuyerSignin()
	SellerSignup()
	SellerSignin()

	fmt.Println()
	fmt.Println("AUTH FLOW COMPLETED")
	fmt.Println()
}

func SellerSignup() {

	req := map[string]any{
		"firstName": "Seller",
		"lastName":  "Test",
		"emailId":   SellerEmail,
		"age":       30,
		"password":  SellerPassword,
		"role":      "seller",
	}

	body, status, err := sellerClient.Post(
		"/api/auth/signup",
		req,
	)

	Assert(err)

	if status != http.StatusCreated {
		panic(string(body))
	}

	AssertStatus(http.StatusCreated, status)

	resp := Decode[AuthResponse](body)

	SellerID = resp.User.ID

	AssertTrue(
		SellerID != "",
		"Seller Signup",
		"seller id is empty",
	)
}

func SellerSignin() {

	req := map[string]any{
		"emailId":  SellerEmail,
		"password": SellerPassword,
	}

	body, status, err := sellerClient.Post(
		"/api/auth/signin",
		req,
	)

	Assert(err)

	if status != http.StatusOK {
		panic(string(body))
	}

	AssertStatus(http.StatusOK, status)

	resp := Decode[AuthResponse](body)

	SellerID = resp.User.ID

	AssertTrue(
		SellerID != "",
		"Seller Signin",
		"seller id is empty",
	)
}

func BuyerSignup() {

	req := map[string]any{
		"firstName": "Buyer",
		"lastName":  "Test",
		"emailId":   BuyerEmail,
		"age":       28,
		"password":  BuyerPassword,
		"role":      "buyer",
	}

	body, status, err := buyerClient.Post(
		"/api/auth/signup",
		req,
	)

	Assert(err)

	if status != http.StatusCreated {
		panic(string(body))
	}

	AssertStatus(http.StatusCreated, status)

	resp := Decode[AuthResponse](body)

	BuyerID = resp.User.ID

	AssertTrue(
		BuyerID != "",
		"Buyer Signup",
		"buyer id is empty",
	)
}

func BuyerSignin() {

	req := map[string]any{
		"emailId":  BuyerEmail,
		"password": BuyerPassword,
	}

	body, status, err := buyerClient.Post(
		"/api/auth/signin",
		req,
	)

	Assert(err)

	if status != http.StatusOK {
		panic(string(body))
	}

	AssertStatus(http.StatusOK, status)

	resp := Decode[AuthResponse](body)

	BuyerID = resp.User.ID

	AssertTrue(
		BuyerID != "",
		"Buyer Signin",
		"buyer id is empty",
	)
}
