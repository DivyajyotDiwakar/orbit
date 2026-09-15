package verify

import (
	"Orbit/internal/utils"

	"github.com/gin-gonic/gin"
)

func VerifySeller(c *gin.Context) {
	var req RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	claim, ok := c.Get("UserFields")
	if !ok {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claim.(*utils.Claims)
	if !ok {
		c.JSON(400, gin.H{"error": "User not authenticated"})
		return
	}
	if userClaims.IsEmailVerified == false {
		c.JSON(403, gin.H{"error": "Email not verified, verify it first.."})
		return
	}

	err := VerifySellerService(userClaims, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to verify seller"})
		return
	}

	c.JSON(200, gin.H{"message": "Seller verification request submitted successfully"})
}
