package google

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.GET("/login", Login())
	r.GET("/callback", Callback())
}
