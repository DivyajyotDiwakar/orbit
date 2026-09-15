package routergroup

import (
	"Orbit/internal/admin"
	"Orbit/internal/auth"
	"Orbit/internal/buyer"
	"Orbit/internal/seller"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(r *gin.RouterGroup) {
	// declare all route groups.
	AuthRouterGroup := r.Group("/auth")
	SellerRouterGroup := r.Group("/seller")
	BuyerRouterGroup := r.Group("/buyer")
	AdminRouterGroup := r.Group("/admin")
	// implement grouping.
	auth.AuthRoutes(AuthRouterGroup)
	seller.SellerRoutes(SellerRouterGroup)
	buyer.BuyerRoutes(BuyerRouterGroup)
	admin.AdminRoutes(AdminRouterGroup)
}
