package buyer

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrSoldOut         = errors.New("sold out")
	ErrAlreadyBooked   = errors.New("already booked")
	ErrProductNotFound = errors.New("product not found")
	ErrSaleNotActive   = errors.New("sale is paused or has ended")
	ErrInvalidEventID  = errors.New("invalid event ID")
	ErrInvalidUserID   = errors.New("invalid user ID")
)

type Service struct {
	requirePayment bool
	paymentWindow  time.Duration
}

func NewService(requirePayment bool, paymentWindow time.Duration) *Service {
	return &Service{requirePayment: requirePayment, paymentWindow: paymentWindow}
}

func (s *Service) Buy(ctx context.Context, productId, eventId string, claim *utils.Claims) (*PurchaseResponse, error) {
	status := StatusConfirmed
	ttl := time.Duration(0)
	if s.requirePayment {
		status = StatusPendingPayment
		ttl = s.paymentWindow
	}

	reservation, result, err := repositories.ReserveProduct(ctx, productId, eventId, claim.ID, string(status), ttl)
	if err != nil {
		return nil, fmt.Errorf("buyer service: %w", err)
	}

	switch result {
	case repositories.BookingSoldOut:
		return nil, ErrSoldOut
	case repositories.BookingAlreadyDone:
		return nil, ErrAlreadyBooked
	case repositories.BookingProductMissing:
		return nil, ErrProductNotFound
	case repositories.BookingSaleNotActive:
		return nil, ErrSaleNotActive
	case repositories.BookingSuccess:
		// go s.persistOrder(reservation, productId, eventId, claim.ID)

		resp := &PurchaseResponse{
			ReservationID: reservation.ReservationID,
			Price:         reservation.Price,
			Status:        reservation.Status,
		}
		if !reservation.ExpiresAt.IsZero() {
			resp.ExpiresAt = &reservation.ExpiresAt
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("buyer service: unexpected result %d", result)
	}
}

// func (s *Service) persistOrder(reservation *repositories.Reservation, productId, eventId, userId string) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	userObjId, err := bson.ObjectIDFromHex(userId)
// 	if err != nil {
// 		log.Printf("WARN persistOrder: invalid userId %s: %v", userId, err)
// 		return
// 	}
// 	productObjId, err := bson.ObjectIDFromHex(productId)
// 	if err != nil {
// 		log.Printf("WARN persistOrder: invalid productId %s: %v", productId, err)
// 		return
// 	}
// 	eventObjId, err := bson.ObjectIDFromHex(eventId)
// 	if err != nil {
// 		log.Printf("WARN persistOrder: invalid eventId %s: %v", eventId, err)
// 		return
// 	}

// 	now := time.Now()
// 	order := models.Order{
// 		ID:            bson.NewObjectID(),
// 		UserID:        userObjId,
// 		ProductID:     productObjId,
// 		EventID:       eventObjId,
// 		ReservationID: reservation.ReservationID,
// 		Price:         reservation.Price,
// 		Status:        models.OrderStatus(reservation.Status),
// 		CreatedAt:     now,
// 		UpdatedAt:     now,
// 	}

// 	if err := repositories.CreateOrder(ctx, order); err != nil {
// 		log.Printf("WARN persistOrder: save failed for reservation=%s: %v", reservation.ReservationID, err)
// 	}
// }

func (s *Service) GetLiveEvents(ctx context.Context) ([]EventView, error) {
	events, err := repositories.GetLiveEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("get live events: %w", err)
	}
	views := make([]EventView, len(events))
	for i, e := range events {
		views[i] = EventView{
			ID:          e.ID.Hex(),
			EventName:   e.EventName,
			Description: e.Description,
			ImageBanner: e.ImageBanner,
			ScheduledAt: e.ScheduledAt,
		}
	}
	return views, nil
}

func (s *Service) GetEventProducts(ctx context.Context, eventIdStr string) ([]ProductView, error) {
	eventId, err := bson.ObjectIDFromHex(eventIdStr)
	if err != nil {
		return nil, ErrInvalidEventID
	}
	products, err := repositories.GetEventProductsWithStock(ctx, eventId)
	if err != nil {
		return nil, fmt.Errorf("get event products: %w", err)
	}
	views := make([]ProductView, len(products))
	for i, p := range products {
		views[i] = ProductView{
			ProductID:      p.ID.Hex(),
			Title:          p.Title,
			Description:    p.Description,
			Price:          p.Price,
			Currency:       p.Currency,
			Image:          p.Image,
			AvailableStock: p.AvailableStock,
		}
	}
	return views, nil
}

func (s *Service) GetMyOrders(ctx context.Context, userIdStr string) ([]OrderView, error) {
	userId, err := bson.ObjectIDFromHex(userIdStr)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	orders, err := repositories.GetOrdersByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get my orders: %w", err)
	}
	if len(orders) == 0 {
		return []OrderView{}, nil
	}

	productIds := make([]bson.ObjectID, 0, len(orders))
	eventIds := make([]bson.ObjectID, 0, len(orders))
	seenProduct := make(map[bson.ObjectID]bool)
	seenEvent := make(map[bson.ObjectID]bool)

	for _, o := range orders {
		if !seenProduct[o.ProductID] {
			productIds = append(productIds, o.ProductID)
			seenProduct[o.ProductID] = true
		}
		if !seenEvent[o.EventID] {
			eventIds = append(eventIds, o.EventID)
			seenEvent[o.EventID] = true
		}
	}

	productTitles, err := repositories.GetProductTitlesByIDs(ctx, productIds)
	if err != nil {
		return nil, fmt.Errorf("get my orders: %w", err)
	}
	eventNames, err := repositories.GetEventNamesByIDs(ctx, eventIds)
	if err != nil {
		return nil, fmt.Errorf("get my orders: %w", err)
	}

	views := make([]OrderView, len(orders))
	for i, o := range orders {
		views[i] = OrderView{
			OrderID:       o.ID.Hex(),
			ProductID:     o.ProductID.Hex(),
			ProductTitle:  productTitles[o.ProductID.Hex()], // "" if product was since deleted
			EventID:       o.EventID.Hex(),
			EventName:     eventNames[o.EventID.Hex()],
			Price:         o.Price,
			Status:        string(o.Status),
			ReservationID: o.ReservationID,
			BookedAt:      o.CreatedAt,
		}
	}
	return views, nil
}
