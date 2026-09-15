package admin

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func extractAdminClaim(c *gin.Context) (*utils.Claims, bool) {
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

func respondVerificationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repositories.ErrVerificationNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "verification request not found"})
	default:
		var stateErr *repositories.VerificationStateError
		if errors.As(err, &stateErr) {
			c.JSON(http.StatusConflict, gin.H{
				"error":         stateErr.Error(),
				"currentStatus": stateErr.Current,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetPendingVerificationsHandler(c *gin.Context) {
	if _, ok := extractAdminClaim(c); !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	reqs, err := GetPendingVerificationsService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pending verifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verifications": reqs})
}

func GetVerificationDetailHandler(c *gin.Context) {
	if _, ok := extractAdminClaim(c); !ok {
		return
	}

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	v, err := GetVerificationDetailService(ctx, id)
	if err != nil {
		respondVerificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"verification": v})
}

func ApproveVerificationHandler(c *gin.Context) {
	claim, ok := extractAdminClaim(c)
	if !ok {
		return
	}

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := ApproveVerificationService(ctx, id, claim.ID); err != nil {
		respondVerificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "seller verification approved"})
}

func RejectVerificationHandler(c *gin.Context) {
	claim, ok := extractAdminClaim(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req RejectRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := RejectVerificationService(ctx, id, claim.ID, req.Reason); err != nil {
		respondVerificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "seller verification rejected"})
}
