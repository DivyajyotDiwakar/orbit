package main

import "fmt"

var (
	sellerClient *Client
	buyerClient  *Client
	adminClient  *Client
)

func main() {

	fmt.Println("========================================")
	fmt.Println("ORBIT BENCHMARK TESTS")
	fmt.Println("========================================")
	fmt.Println()

	sellerClient = NewClient()
	buyerClient = NewClient()
	adminClient = NewClient()

	fmt.Println("========== AUTH ==========")

	AuthFlow()

	// AdminFlow() 

	SellerSignin()

	fmt.Println()
	fmt.Println("AUTH COMPLETED")
	fmt.Println()
	
	SellerFlow()

	fmt.Println()
	fmt.Println("========== BENCHMARK REPORT ==========")
	PrintBenchmark()

	PrintSummary()
}
