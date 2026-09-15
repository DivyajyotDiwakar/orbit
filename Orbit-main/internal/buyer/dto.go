package buyer

import "time"

type BookingStatus string

const (
	StatusConfirmed      BookingStatus = "CONFIRMED"
	StatusPendingPayment BookingStatus = "PENDING_PAYMENT"
)

type PurchaseResponse struct {
	ReservationID string     `json:"reservationId"`
	Price         float64    `json:"price"`
	Status        string     `json:"status"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
}

type EventView struct {
	ID          string    `json:"id"`
	EventName   string    `json:"eventName"`
	Description string    `json:"description"`
	ImageBanner string    `json:"imageBanner"`
	ScheduledAt time.Time `json:"scheduledAt"`
}

type ProductView struct {
	ProductID      string  `json:"productId"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Currency       string  `json:"currency"`
	Image          string  `json:"image"`
	AvailableStock int     `json:"availableStock"`
}

type OrderView struct {
	OrderID       string    `json:"orderId"`
	ProductID     string    `json:"productId"`
	ProductTitle  string    `json:"productTitle"`
	EventID       string    `json:"eventId"`
	EventName     string    `json:"eventName"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	ReservationID string    `json:"reservationId"`
	BookedAt      time.Time `json:"bookedAt"`
}
