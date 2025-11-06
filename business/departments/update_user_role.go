package departments

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][update_user_role]")

	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			code := response.ErrorResponse("User ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req UpdateUserRoleRequest
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
			code := response.ErrorResponse("Only leaders can update user roles")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Validate role
		var newRole usercol.Role
		switch req.Role {
		case "employee":
			newRole = usercol.RoleEmployee
		case "manager":
			newRole = usercol.RoleManager
		case "leader":
			newRole = usercol.RoleLeader
		case "assistant_director":
			newRole = usercol.RoleAssistantDirector
		default:
			code := response.ErrorResponse("Invalid role. Must be one of: employee, manager, leader, assistant_director")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Kiểm tra user có tồn tại không
		targetUser, err := usercol.FindWithUserID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("User not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get user")
			code := response.ErrorResponse("Failed to get user")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Không cho phép leader tự đổi role của chính mình
		if targetUser.GetIDString() == user.GetIDString() {
			code := response.ErrorResponse("You cannot change your own role")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Cập nhật role
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			code := response.ErrorResponse("Invalid user ID")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		updateData := bson.M{
			"role": newRole,
		}

		updatedUser, err := usercol.UpdateByID(c.Request.Context(), objID, updateData)
		if err != nil {
			logger.Err(err).Msg("failed to update user role")
			code := response.ErrorResponse("Failed to update user role")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Format response
		responseData := map[string]interface{}{
			"id":        updatedUser.GetIDString(),
			"full_name": updatedUser.FullName,
			"email":     updatedUser.Email,
			"avatar":    updatedUser.Avatar,
			"role":      updatedUser.Role,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

