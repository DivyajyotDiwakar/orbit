package order

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func extractClaim(c *gin.Context) (*utils.Claims, bool) {
	raw, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return nil, false
	}
	claim, ok := raw.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return nil, false
	}
	return claim, true
}

func GetSellerEventOrdersHandler(c *gin.Context) {
	eventIdStr := c.Param("eventId")

	eventId, err := bson.ObjectIDFromHex(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orders, err := repositories.GetOrdersByEvent(ctx, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func GetEventAnalyticsHandler(c *gin.Context) {
	eventIdStr := c.Param("eventId")

	eventId, err := bson.ObjectIDFromHex(eventIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	analytics, err := repositories.GetEventAnalytics(ctx, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

func GetSellerAllOrdersHandler(c *gin.Context) {
	claim, ok := extractClaim(c)
	if !ok {
		return
	}

	sellerId, err := bson.ObjectIDFromHex(claim.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	orders, err := repositories.GetOrdersBySeller(ctx, sellerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"orders": []SellerOrderView{}})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve product info"})
		return
	}
	eventNames, err := repositories.GetEventNamesByIDs(ctx, eventIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event info"})
		return
	}

	views := make([]SellerOrderView, len(orders))
	for i, o := range orders {
		views[i] = SellerOrderView{
			OrderID:       o.ID.Hex(),
			ProductID:     o.ProductID.Hex(),
			ProductTitle:  productTitles[o.ProductID.Hex()],
			EventID:       o.EventID.Hex(),
			EventName:     eventNames[o.EventID.Hex()],
			BuyerID:       o.UserID.Hex(),
			Price:         o.Price,
			Status:        string(o.Status),
			ReservationID: o.ReservationID,
			BookedAt:      o.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"orders": views})
}
