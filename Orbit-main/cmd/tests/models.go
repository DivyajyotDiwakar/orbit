package main

import "time"

const TestProductImage = "cmd/tests/assets/product.jpg"

var (
	EventID        string
	ProductID      string
	VerificationID string
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateEventRequest struct {
	Title       string
	Description string
	ScheduledAt string
}

type UpdateEventRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ScheduledAt string `json:"scheduledAt,omitempty"`
	ImageBanner string `json:"imageBanner,omitempty"`
}

type CreateEventResponse struct {
	Message string `json:"message"`
	EventID string `json:"eventId"`
}

type SellerEvent struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"sellerId"`
	EventName   string    `json:"eventName"`
	Description string    `json:"description"`
	ScheduledAt time.Time `json:"scheduledAt"`
	IsLive      bool      `json:"isLive"`
	ImageBanner string    `json:"imageBanner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EventResponse struct {
	Event SellerEvent `json:"event"`
}

type EventsResponse struct {
	Events []SellerEvent `json:"events"`
}

type Inventory struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"sellerId"`
	EventID     string    `json:"eventId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Frequency   int       `json:"frequency"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProductResponse struct {
	Product Inventory `json:"product"`
}

type ProductsResponse struct {
	Products []Inventory `json:"products"`
}

type LiveEvent struct {
	ID          string    `json:"id"`
	EventName   string    `json:"eventName"`
	Description string    `json:"description"`
	ImageBanner string    `json:"imageBanner"`
	ScheduledAt time.Time `json:"scheduledAt,omitempty"`
}

type LiveEventsResponse struct {
	Events []LiveEvent `json:"events"`
}

type BuyerProduct struct {
	ProductID      string  `json:"productId"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Currency       string  `json:"currency"`
	Image          string  `json:"image"`
	AvailableStock int     `json:"availableStock"`
}

type BuyerProductsResponse struct {
	Products []BuyerProduct `json:"products"`
}

type PurchaseResponse struct {
	ReservationID string  `json:"reservationId"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
}

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ProductID string    `json:"productId"`
	EventID   string    `json:"eventId"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	BookedAt  time.Time `json:"bookedAt"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

type AnalyticsResponse struct {
	Analytics any `json:"analytics"`
}

type PendingVerification struct {
	ID string `json:"id"`
}

type PendingVerificationResponse struct {
	Verifications []PendingVerification `json:"verifications"`
}

type VerificationDetail struct {
	ID string `json:"id"`
}

type VerificationDetailResponse struct {
	Verification VerificationDetail `json:"verification"`
}

type RejectVerificationRequest struct {
	Reason string `json:"reason"`
}
