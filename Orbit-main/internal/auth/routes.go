package auth

import (
	"Orbit/internal/auth/check"
	"Orbit/internal/auth/signin"
	"Orbit/internal/auth/signup"
	"Orbit/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/signin", signin.SignInHandler)
	r.POST("/signup", signup.Register)
	r.GET("/check", middleware.UserMiddleware(), check.CheckHandler)
}
