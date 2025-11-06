package departments

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RemoveEmployee() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][remove_employee]")

	return func(c *gin.Context) {
		id := c.Param("id")
		employeeID := c.Param("employee_id")

		if id == "" {
			code := response.ErrorResponse("Team ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		if employeeID == "" {
			code := response.ErrorResponse("Employee ID is required")
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
			code := response.ErrorResponse("Only leaders can remove employees from teams")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Tìm team
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

		// Kiểm tra team có manager chưa
		if team.ManagerID == "" {
			code := response.ErrorResponse("Team does not have a manager")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Kiểm tra employee có thuộc team này không
		belongs, err := teammembercol.CheckEmployeeBelongsToManager(c.Request.Context(), team.ManagerID, employeeID)
		if err != nil {
			logger.Err(err).Msg("failed to check employee relationship")
			code := response.ErrorResponse("Failed to check employee relationship")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		if !belongs {
			code := response.ErrorResponse("Employee is not in this team")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Xóa team member relationship
		err = teammembercol.SoftDeleteByManagerAndEmployee(c.Request.Context(), team.ManagerID, employeeID)
		if err != nil {
			logger.Err(err).Msg("failed to remove employee from team")
			code := response.ErrorResponse("Failed to remove employee from team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{
			"message":     "Employee removed from team successfully",
			"team_id":     team.GetIDString(),
			"employee_id": employeeID,
		}))
	}
}

