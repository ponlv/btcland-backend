package routers

import (
	"api/business/auth"
	"api/business/healthcheck"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("healthcheck", healthcheck.Healthcheck())

	// Auth routes
	authRouter := r.Group("auth")
	auth.AddRouter(authRouter)

}
