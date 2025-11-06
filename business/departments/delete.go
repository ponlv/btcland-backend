package departments

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Delete() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][delete]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Lấy user từ context
		userInterface, exists := c.Get("current_user")
		if !exists {
			code := response.ErrorResponse("Unauthorized")
			c.JSON(http.StatusUnauthorized, code)
			c.Abort()
			return
		}

		user, ok := userInterface.(*usercol.User)
		if !ok {
			code := response.ErrorResponse("Invalid user")
			c.JSON(http.StatusUnauthorized, code)
			c.Abort()
			return
		}

		// Kiểm tra role - chỉ leader mới có quyền
		if user.Role != usercol.RoleLeader {
			code := response.ErrorResponse("Only leaders can delete teams")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Kiểm tra team có tồn tại không
		team, err := teamcol.FindByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Team not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get team")
			code := response.ErrorResponse("Failed to get team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Xóa mềm team
		err = teamcol.SoftDelete(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Msg("failed to delete team")
			code := response.ErrorResponse("Failed to delete team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{
			"message": "Team deleted successfully",
			"id":      team.GetIDString(),
		}))
	}
}

