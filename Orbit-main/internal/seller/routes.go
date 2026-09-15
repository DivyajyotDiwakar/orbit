package seller

import (
	"Orbit/internal/middleware"
	"Orbit/internal/order"
	"Orbit/internal/seller/event"
	"Orbit/internal/seller/isVerified"
	"Orbit/internal/seller/verify"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(r *gin.RouterGroup) {
	r.Use(
		middleware.UserMiddleware(),
		middleware.SellerMiddleware(),
	)

	r.GET("/isVerified", isVerified.IsVerifiedSeller)
	r.POST("/verify", verify.VerifySeller)

	r.GET("/orders", order.GetSellerEventOrdersHandler)

	EventGroups := r.Group("/events")
	event.EventRoutes(EventGroups)
}
