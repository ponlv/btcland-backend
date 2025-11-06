package healthcheck

import (
	"api/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthcheck() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.JSON(http.StatusOK, response.SuccessResponse(nil))
		return
	}
}
