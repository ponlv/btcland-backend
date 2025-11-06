package images

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	// Public route - không cần authentication để xem ảnh
	r.GET("*path", Get()) // GET /images/*path
}

