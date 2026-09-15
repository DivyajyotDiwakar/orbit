package isVerified

import (
	"Orbit/internal/repositories"
	"Orbit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsVerifiedSeller(c *gin.Context) {
	value, exists := c.Get("UserFields")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	claim, ok := value.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Not Authenticated",
		})
		return
	}

	objId, err := bson.ObjectIDFromHex(claim.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := repositories.FindUserById(c.Request.Context(), objId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	isApproved := user.SellerInfo != nil && user.SellerInfo.IsApproved

	c.JSON(http.StatusOK, gin.H{
		"isApproved": isApproved,
	})
}
