package main

var (
	SellerToken string
	BuyerToken  string
	AdminToken  string

	ProductIDs []string

	ProductFrequency     = map[string]int{}
	SellerCapturedOrders []Order
)
