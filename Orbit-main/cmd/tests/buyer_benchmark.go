package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BuyerBenchmark() {
	GetLiveEvents()
	GetEventProducts()
	// createBuyers()

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 1; i <= BuyerCount; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			<-start
			BuyerWorker(id)
		}(i)
	}

	BenchmarkStart()
	close(start)
	wg.Wait()

	BenchmarkFinish()

	completed := Metrics.Success + Metrics.SoldOut + Metrics.AlreadyBooked + Metrics.Errors
	AssertEqual(int64(BuyerCount), completed, "Benchmark Request Accounting")
}

func createBuyers() {
	type job struct{ index int }

	jobs := make(chan job, BuyerCount)
	var wg sync.WaitGroup

	for w := 0; w < WorkerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				client := NewClient()
				email := BuyerPrefix + strconv.Itoa(j.index) + "@gmail.com"
				SignupBuyer(client, email, BuyerPassword)
			}
		}()
	}

	for i := 1; i <= BuyerCount; i++ {
		jobs <- job{index: i}
	}
	close(jobs)

	wg.Wait()
}

func BuyerWorker(_ int) {
	productID := RandomProduct()
	BuyProduct(buyerClient, productID, bson.NewObjectID().Hex())
}

func GetLiveEvents() {

	body, status, err := buyerClient.Get(
		"/api/buyer/events",
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[LiveEventsResponse](body)

	AssertTrue(
		len(resp.Events) > 0,
		"Get Live Events",
		"no live events found",
	)
}

func GetEventProducts() {

	body, status, err := buyerClient.Get(
		"/api/buyer/event/" + EventID,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[BuyerProductsResponse](body)

	AssertTrue(
		len(resp.Products) > 0,
		"Get Event Products",
		"no products returned",
	)

	ProductIDs = ProductIDs[:0]

	for _, p := range resp.Products {
		ProductIDs = append(ProductIDs, p.ProductID)
	}
}

func SignupBuyer(client *Client, email, password string) {

	req := map[string]any{
		"firstName": "Benchmark",
		"lastName":  "Buyer",
		"emailId":   email,
		"age":       25,
		"password":  password,
		"role":      "buyer",
	}

	body, status, err := client.Post(
		"/api/auth/signup",
		req,
	)

	if err != nil {
		AddErrorSample("request failed: " + err.Error())
		IncError()
		return
	}

	if status != http.StatusCreated {
		IncError()
		return
	}

	resp := Decode[AuthResponse](body)

	if resp.User.ID == "" {
		IncError()
	}
}

func SigninBuyer(client *Client, email, password string) {

	req := map[string]any{
		"emailId":  email,
		"password": password,
	}

	body, status, err := client.Post(
		"/api/auth/signin",
		req,
	)

	if err != nil {
		IncError()
		return
	}

	if status != http.StatusOK {
		IncError()
		return
	}

	resp := Decode[AuthResponse](body)

	if resp.User.ID == "" {
		IncError()
	}
}

func RandomProduct() string {
	return ProductIDs[rand.Intn(len(ProductIDs))]
}

func BuyProduct(client *Client, productID string, userID string) {

	start := time.Now()

	path := fmt.Sprintf(
		"/api/buyer/event/%s/purchase/%s?userId=%s",
		EventID,
		productID,
		userID,
	)

	body, status, err := client.Post(path, nil)

	AddLatency(time.Since(start))

	if err != nil {
		AddErrorSample("purchase request failed: " + err.Error())
		IncError()
		return
	}

	switch status {

	case http.StatusCreated:
		resp := Decode[PurchaseResponse](body)

		if resp.Status != "" {
			IncSuccess()
		} else {
			IncError()
		}

	case http.StatusConflict:
		resp := Decode[ErrorResponse](body)

		switch resp.Error {
		case "sold out":
			IncSoldOut()

		case "already booked":
			IncAlreadyBooked()

		default:
			AddErrorSample("status=409 body=" + compactResponse(body))
			IncError()
		}

	default:
		AddErrorSample(fmt.Sprintf("status=%d body=%s", status, compactResponse(body)))
		IncError()
	}
}

func compactResponse(body []byte) string {
	const maxLength = 300
	response := strings.TrimSpace(string(body))
	if len(response) > maxLength {
		return response[:maxLength] + "..."
	}
	return response
}
