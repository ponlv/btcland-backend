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

type AssignManagerRequest struct {
	ManagerID string `json:"manager_id" binding:"required"`
}

func AssignManager() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][assign_manager]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req AssignManagerRequest
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
			code := response.ErrorResponse("Only leaders can assign managers")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Kiểm tra team có tồn tại không
		_, err := teamcol.FindByID(c.Request.Context(), id)
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

		// Kiểm tra manager có tồn tại và có role là manager không
		manager, err := usercol.FindWithUserID(c.Request.Context(), req.ManagerID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Manager not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get manager")
			code := response.ErrorResponse("Failed to get manager")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		if manager.Role != usercol.RoleManager {
			code := response.ErrorResponse("User must be a manager")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Cập nhật manager_id cho team
		err = teamcol.UpdateManagerID(c.Request.Context(), id, req.ManagerID)
		if err != nil {
			logger.Err(err).Msg("failed to assign manager")
			code := response.ErrorResponse("Failed to assign manager")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại team đã cập nhật
		updated, err := teamcol.FindByID(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Msg("failed to get updated team")
			code := response.ErrorResponse("Manager assigned but failed to retrieve team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy thông tin manager
		managerData := map[string]interface{}{
			"id":        manager.GetIDString(),
			"full_name": manager.FullName,
			"email":     manager.Email,
			"avatar":    manager.Avatar,
		}

		responseData := map[string]interface{}{
			"id":          updated.GetIDString(),
			"name":        updated.Name,
			"description": updated.Description,
			"manager_id":  updated.ManagerID,
			"manager":     managerData,
			"created_at":  updated.CreatedAt,
			"updated_at":  updated.UpdatedAt,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

