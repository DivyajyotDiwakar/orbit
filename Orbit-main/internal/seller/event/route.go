package event

import (
	"Orbit/internal/order"
	"Orbit/internal/seller/event/product"
	"Orbit/internal/seller/event/sale"

	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.RouterGroup) {
	r.POST("/create", CreateAnEventHandler)
	r.GET("/get", GetAllEventsHandler)
	r.GET("/get/:id", GetAnEventHandler)
	r.PUT("/update/:id", UpdateAnEventHandler)
	r.DELETE("/delete/:id", DeleteAnEventHandler)

	r.POST("/:eventId/registerProducts", product.RegisterProductHandler)
	r.GET("/:eventId/getProducts", product.GetAllEventProductsHandler)
	r.GET("/:eventId/getProduct/:productId", product.GetAnEventProductHandler)
	r.PUT("/:eventId/updateProduct/:productId", product.UpdateAnEventProductHandler)
	r.DELETE("/:eventId/deleteProduct/:productId", product.DeleteAnEventProductHandler)

	r.POST("/:eventId/Live", sale.LiveSaleHandler)
	r.POST("/:eventId/Pause", sale.PauseSaleHandler)
	r.POST("/:eventId/Resume", sale.ResumeSaleHandler)
	r.POST("/:eventId/End", sale.StopSaleHandler)

	r.GET("/booking/fetch/:eventId", order.GetSellerEventOrdersHandler)
	r.GET("/booking/analytics/:eventId", order.GetEventAnalyticsHandler)
}
