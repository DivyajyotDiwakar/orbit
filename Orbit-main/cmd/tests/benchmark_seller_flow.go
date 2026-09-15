package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func SellerFlow() {

	fmt.Println("========== SELLER FLOW ==========")

	CreateEvent()

	GetMyEvents()

	GetEvent()

	UpdateEvent()

	RegisterProducts()

	GetProducts()

	UpdateProducts()

	LiveSale()

	BuyerBenchmark()

	SellerSignin()

	PauseSale()

	ResumeSale()

	EndSale()

	GetSellerOrders()

	VerifyInventoryConsistency()

	GetEventAnalytics()

	DeleteProducts()

	DeleteEvent()

	fmt.Println()
	fmt.Println("SELLER FLOW COMPLETED")
	fmt.Println()
}

func CreateEvent() {

	req := CreateEventRequest{
		Title:       "Orbit Benchmark Event",
		Description: "Created by benchmark tests",
		ScheduledAt: time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02T15:04"),
	}

	body, status, err := sellerClient.UploadMultipart(
		http.MethodPost,
		"/api/seller/events/create",
		map[string]string{
			"title":       req.Title,
			"description": req.Description,
			"scheduledAt": req.ScheduledAt,
		},
		"imageBanner",
		TestProductImage,
	)

	Assert(err)
	fmt.Println("STATUS:", status)
	fmt.Println("BODY:", string(body))

	AssertStatus(http.StatusCreated, status)

	resp := Decode[CreateEventResponse](body)

	EventID = resp.EventID

	AssertTrue(
		EventID != "",
		"Create Event",
		"event id is empty",
	)
}

func GetMyEvents() {

	body, status, err := sellerClient.Get(
		"/api/seller/events/get",
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[EventsResponse](body)

	AssertTrue(
		len(resp.Events) > 0,
		"Get Events",
		"seller has no events",
	)
}

func GetEvent() {

	body, status, err := sellerClient.Get(
		"/api/seller/events/get/" + EventID,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[EventResponse](body)

	AssertTrue(
		resp.Event.ID == EventID,
		"Get Event",
		"invalid event returned",
	)
}

func UpdateEvent() {

	req := UpdateEventRequest{
		Title:       "Updated Orbit Benchmark Event",
		Description: "Updated by benchmark tests",
		ScheduledAt: time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339),
		ImageBanner: "updated-banner",
	}

	body, status, err := sellerClient.Put(
		"/api/seller/events/update/"+EventID,
		req,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "event updated successfully",
		"Update Event",
		resp.Message,
	)
}

func RegisterProducts() {
	const workers = 5

	type job struct {
		index int
	}

	jobs := make(chan job, ProductCount)
	errCh := make(chan error, ProductCount)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := range jobs {
				fields := map[string]string{
					"title":       fmt.Sprintf("%s %d", ProductPrefix, j.index+1),
					"description": "Created by benchmark tests",
					"price":       "249999",
					"frequency":   "1",
				}

				body, status, err := sellerClient.UploadMultipart(
					http.MethodPost,
					"/api/seller/events/"+EventID+"/registerProducts",
					fields,
					"image",
					TestProductImage,
				)

				if err != nil {
					errCh <- err
					continue
				}

				if status != http.StatusCreated {
					errCh <- fmt.Errorf("status=%d body=%s", status, string(body))
				}
			}
		}()
	}

	for i := 0; i < ProductCount; i++ {
		jobs <- job{index: i}
	}
	close(jobs)

	wg.Wait()
	close(errCh)

	for err := range errCh {
		Assert(err)
	}
}

func GetProducts() {

	body, status, err := sellerClient.Get(
		"/api/seller/events/" + EventID + "/getProducts",
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[ProductsResponse](body)

	AssertEqual(
		ProductCount,
		len(resp.Products),
		"Get Products",
	)

	ProductIDs = ProductIDs[:0]

	for _, product := range resp.Products {

		AssertTrue(
			product.ID != "",
			"Capture Product ID",
			"empty product id",
		)

		ProductIDs = append(ProductIDs, product.ID)
		ProductFrequency[product.ID] = product.Frequency
	}

	ProductID = ProductIDs[0]
}

func UpdateProducts() {
	const workers = 5

	type job struct {
		index     int
		productID string
	}

	jobs := make(chan job, len(ProductIDs))
	errCh := make(chan error, len(ProductIDs))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := range jobs {

				fields := map[string]string{
					"title":       fmt.Sprintf("%s %d Updated", ProductPrefix, j.index+1),
					"description": "Updated by benchmark tests",
					"price":       "199999",
					"frequency":   "1",
				}

				body, status, err := sellerClient.UploadMultipart(
					http.MethodPut,
					"/api/seller/events/"+EventID+"/updateProduct/"+j.productID,
					fields,
					"image",
					TestProductImage,
				)

				if err != nil {
					errCh <- err
					continue
				}

				if status != http.StatusOK {
					errCh <- fmt.Errorf("update product %d failed: status=%d body=%s",
						j.index+1,
						status,
						string(body),
					)
					continue
				}

				resp := Decode[MessageResponse](body)

				if resp.Message != "product updated successfully" {
					errCh <- fmt.Errorf("update product %d: %s",
						j.index+1,
						resp.Message,
					)
				}
			}
		}()
	}

	for i, productID := range ProductIDs {
		jobs <- job{
			index:     i,
			productID: productID,
		}
	}

	close(jobs)

	wg.Wait()
	close(errCh)

	for err := range errCh {
		Assert(err)
	}
}

func DeleteProducts() {

	for i, productID := range ProductIDs {

		body, status, err := sellerClient.Delete(
			"/api/seller/events/" + EventID + "/deleteProduct/" + productID,
		)

		Assert(err)
		AssertStatus(http.StatusOK, status)

		resp := Decode[MessageResponse](body)

		AssertTrue(
			resp.Message == "product deleted successfully",
			fmt.Sprintf("Delete Product %d", i+1),
			resp.Message,
		)
	}
}

func LiveSale() {

	body, status, err := sellerClient.Post(
		"/api/seller/events/"+EventID+"/Live",
		nil,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "sale is now live",
		"Live Sale",
		resp.Message,
	)
}

func PauseSale() {

	body, status, err := sellerClient.Post(
		"/api/seller/events/"+EventID+"/Pause",
		nil,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "sale paused — no new bookings will be accepted",
		"Pause Sale",
		resp.Message,
	)
}

func ResumeSale() {

	body, status, err := sellerClient.Post(
		"/api/seller/events/"+EventID+"/Resume",
		nil,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "sale resumed",
		"Resume Sale",
		resp.Message,
	)
}

func EndSale() {

	body, status, err := sellerClient.Post(
		"/api/seller/events/"+EventID+"/End",
		nil,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "sale ended — all orders synced to database",
		"End Sale",
		resp.Message,
	)
}

func GetSellerOrders() {

	body, status, err := sellerClient.Get(
		"/api/seller/events/booking/fetch/" + EventID,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[OrdersResponse](body)

	AssertTrue(
		resp.Orders != nil,
		"Get Seller Orders",
		"orders response is nil",
	)

	SellerCapturedOrders = resp.Orders
}

func VerifyInventoryConsistency() {

	soldCount := map[string]int{}
	for _, o := range SellerCapturedOrders {
		soldCount[o.ProductID]++
	}

	oversold := false
	for productID, sold := range soldCount {
		stock, known := ProductFrequency[productID]
		if !known {
			continue // product not tracked (shouldn't happen, but don't false-positive)
		}
		if sold > stock {
			oversold = true
			fmt.Printf("   ↳ OVERSOLD product=%s sold=%d stock=%d\n", productID, sold, stock)
		}
	}

	AssertTrue(
		!oversold,
		"Zero Overselling",
		"one or more products sold more units than their registered stock",
	)

	successCount := int(Metrics.Success)
	mongoOrderCount := len(SellerCapturedOrders)

	AssertEqual(
		successCount,
		mongoOrderCount,
		"Inventory Consistency (Redis success count == MongoDB order count)",
	)
}

func GetEventAnalytics() {

	body, status, err := sellerClient.Get(
		"/api/seller/events/booking/analytics/" + EventID,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[AnalyticsResponse](body)

	AssertTrue(
		resp.Analytics != nil,
		"Get Event Analytics",
		"analytics response is nil",
	)
}

func DeleteEvent() {

	body, status, err := sellerClient.Delete(
		"/api/seller/events/delete/" + EventID,
	)

	Assert(err)
	AssertStatus(http.StatusOK, status)

	resp := Decode[MessageResponse](body)

	AssertTrue(
		resp.Message == "event deleted successfully",
		"Delete Event",
		resp.Message,
	)
}
