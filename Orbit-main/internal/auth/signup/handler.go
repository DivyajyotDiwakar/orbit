package signup

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 7*time.Second)
	defer cancel()

	tokenString, user, err := RegisterUserService(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", tokenString, 7200, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Register Successfully",
		"user": RegisterResponse{
			ID:        user.ID.Hex(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			EmailId:   user.EmailId,
			Age:       user.Age,
			Role:      user.Role,
		},
	})
}
