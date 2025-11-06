package common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetToken(c *gin.Context) string {
	auth := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(auth, "Bearer ")

	return token
}