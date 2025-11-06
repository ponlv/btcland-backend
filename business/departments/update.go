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

type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Update() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][update]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req UpdateRequest
		if err := c.BindJSON(&req); err != nil {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
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
			code := response.ErrorResponse("Only leaders can update teams")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Tìm team hiện tại
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

		// Cập nhật thông tin
		if req.Name != "" {
			team.Name = req.Name
		}
		if req.Description != "" {
			team.Description = req.Description
		}

		// Lưu cập nhật
		_, err = teamcol.Update(c.Request.Context(), team)
		if err != nil {
			logger.Err(err).Msg("failed to update team")
			code := response.ErrorResponse("Failed to update team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại team đã cập nhật
		updated, err := teamcol.FindByID(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Msg("failed to get updated team")
			code := response.ErrorResponse("Failed to retrieve updated team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(updated))
	}
}

