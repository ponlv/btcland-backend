package departments

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ManagerID   string `json:"manager_id"` // Optional, có thể gán sau
}

func Create() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][create]")

	return func(c *gin.Context) {
		var req CreateRequest
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
			code := response.ErrorResponse("Only leaders can create teams")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Kiểm tra manager_id nếu có
		if req.ManagerID != "" {
			manager, err := usercol.FindWithUserID(c.Request.Context(), req.ManagerID)
			if err != nil {
				code := response.ErrorResponse("Manager not found")
				c.JSON(http.StatusBadRequest, code)
				c.Abort()
				return
			}

			// Kiểm tra manager có role là manager không
			if manager.Role != usercol.RoleManager {
				code := response.ErrorResponse("User must be a manager")
				c.JSON(http.StatusBadRequest, code)
				c.Abort()
				return
			}
		}

		// Tạo team mới
		team := &teamcol.Team{
			Name:        req.Name,
			Description: req.Description,
			ManagerID:   req.ManagerID,
		}

		_, err := teamcol.Create(c.Request.Context(), team)
		if err != nil {
			logger.Err(err).Msg("failed to create team")
			code := response.ErrorResponse("Failed to create team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại team vừa tạo
		created, err := teamcol.FindByID(c.Request.Context(), team.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get created team")
			code := response.ErrorResponse("Team created but failed to retrieve")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(created))
	}
}
