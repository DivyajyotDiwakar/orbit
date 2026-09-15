package sale

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LiveSaleHandler(c *gin.Context) {
	eventId := c.Param("eventId")
	claim, ok := extractClaim(c)
	if !ok {
		return
	}
	if err := LiveSaleService(claim.ID, eventId); err != nil {
		respondSaleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sale is now live"})
}

func PauseSaleHandler(c *gin.Context) {
	eventId := c.Param("eventId")
	claim, ok := extractClaim(c)
	if !ok {
		return
	}
	if err := PauseSaleService(claim.ID, eventId); err != nil {
		respondSaleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sale paused — no new bookings will be accepted"})
}

func ResumeSaleHandler(c *gin.Context) {
	eventId := c.Param("eventId")
	claim, ok := extractClaim(c)
	if !ok {
		return
	}
	if err := ResumeSaleService(claim.ID, eventId); err != nil {
		respondSaleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sale resumed"})
}

func StopSaleHandler(c *gin.Context) {
	eventId := c.Param("eventId")
	claim, ok := extractClaim(c)
	if !ok {
		return
	}
	if err := StopSaleService(claim.ID, eventId); err != nil {
		respondSaleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sale ended — all orders synced to database"})
}

func extractClaim(c *gin.Context) (*utils.Claims, bool) {
	raw, ok := c.Get("UserFields")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return nil, false
	}
	claim, ok := raw.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid session"})
		return nil, false
	}
	return claim, true
}

func respondSaleError(c *gin.Context, err error) {
	var transErr *repositories.SaleTransitionError
	if errors.As(err, &transErr) {
		c.JSON(http.StatusConflict, gin.H{
			"error":           transErr.Error(),
			"currentStatus":   string(transErr.Current),
			"attemptedAction": string(transErr.Attempted),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
