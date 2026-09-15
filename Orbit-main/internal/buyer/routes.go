package buyer

import (
	"Orbit/internal/middleware"

	"github.com/gin-gonic/gin"
)

func BuyerRoutes(r *gin.RouterGroup) {
	svc := NewService(false, 0)
	h := NewHandler(svc)

	r.Use(middleware.UserMiddleware())
	r.GET("/events", h.GetLiveEventsHandler)
	r.GET("/event/:id", h.GetEventProductsHandler)

	r.Use(middleware.BuyerMiddleware())

	r.POST("/event/:eventId/purchase/:productId", h.BuyerEventHandler)
	r.GET("/orders", h.GetMyOrdersHandler)
}
