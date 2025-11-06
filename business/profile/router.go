package profile

import (
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Tất cả routes đều yêu cầu authentication
	r.Use(middleware.AuthMiddleware())

	// Routes cho Profile cá nhân
	r.GET("", Get())                 // GET /profile - Lấy thông tin profile
	r.PUT("", Update())              // PUT /profile - Cập nhật thông tin profile
	r.POST("avatar", UpdateAvatar()) // POST /profile/avatar - Cập nhật avatar
}
