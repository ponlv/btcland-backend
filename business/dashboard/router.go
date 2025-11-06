package dashboard

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Tất cả routes đều yêu cầu authentication
	r.Use(middleware.AuthMiddleware())

	// Routes cho Lãnh đạo (Leader) - Dashboard và thống kê
	r.GET("stats", Stats())                                     // GET /dashboard/stats - Thống kê tổng quan
	r.GET("work-confirmations-stats", WorkConfirmationsStats()) // GET /dashboard/work-confirmations-stats - Thống kê đơn xác nhận
	r.GET("teams-stats", TeamsStats())                          // GET /dashboard/teams-stats - Thống kê theo phòng ban/team
}
