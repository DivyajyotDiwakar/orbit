package signin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SignInHandler(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 7*time.Second)
	defer cancel()

	tokenString, user, err := SignInService(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", tokenString, 7200, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Sign In Successfully",
		"user": SignInResponse{
			ID:        user.ID.Hex(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			EmailId:   user.EmailId,
			Age:       user.Age,
			Role:      user.Role,
		},
	})
}
