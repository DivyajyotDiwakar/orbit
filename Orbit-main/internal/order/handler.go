package order

import (
	"Orbit/internal/repositories"
	"Orbit/internal/seller/order"
	"Orbit/internal/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetMyOrdersHandler(c *gin.Context) {
	raw, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not permitted"})
		return
	}
	claim, ok := raw.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return
	}

	userId, err := bson.ObjectIDFromHex(claim.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orders, err := repositories.GetOrdersByUser(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
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

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"orders": []order.SellerOrderView{}})
		return
	}

	userIDs := make([]bson.ObjectID, 0)
	productIDs := make([]bson.ObjectID, 0)
	seenUserIDs := make(map[bson.ObjectID]bool)
	seenProductIDs := make(map[bson.ObjectID]bool)

	for _, o := range orders {
		if !seenUserIDs[o.UserID] {
			userIDs = append(userIDs, o.UserID)
			seenUserIDs[o.UserID] = true
		}
		if !seenProductIDs[o.ProductID] {
			productIDs = append(productIDs, o.ProductID)
			seenProductIDs[o.ProductID] = true
		}
	}

	users, err := repositories.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user data"})
		return
	}

	productTitles, err := repositories.GetProductTitlesByIDs(ctx, productIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve product data"})
		return
	}

	views := make([]order.SellerOrderView, len(orders))
	for i, o := range orders {
		customer := users[o.UserID.Hex()]
		productTitle, ok := productTitles[o.ProductID.Hex()]
		if !ok {
			productTitle = "Unknown Product"
		}
		views[i] = order.SellerOrderView{
			OrderID:       o.ID.Hex(),
			ProductID:     o.ProductID.Hex(),
			ProductTitle:  productTitle,
			CustomerName:  fmt.Sprintf("%s %s", customer.FirstName, customer.LastName),
			CustomerEmail: customer.EmailId,
			Price:         o.Price,
			Status:        string(o.Status),
			BookedAt:      o.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"orders": views})
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
