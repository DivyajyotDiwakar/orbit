package admin

import (
	"Orbit/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	r.Use(
		middleware.UserMiddleware(),
		middleware.AdminMiddleware(),
	)

	r.GET("/verifications/pending", GetPendingVerificationsHandler)
	r.GET("/verifications/:id", GetVerificationDetailHandler)
	r.POST("/verifications/:id/approve", ApproveVerificationHandler)
	r.POST("/verifications/:id/reject", RejectVerificationHandler)
}
