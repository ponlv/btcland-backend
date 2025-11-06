package auth

import (
	"api/business/auth/google"

	"github.com/gin-gonic/gin"
)

func AddRouter(r *gin.RouterGroup) {

	googleGroup := r.Group("google")
	google.Router(googleGroup)

	// QR login - removed
}
