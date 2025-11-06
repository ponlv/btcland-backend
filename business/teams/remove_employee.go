package teams

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

func RemoveEmployee() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][remove_employee]")

	return func(c *gin.Context) {
		employeeID := c.Param("employee_id")
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

		// Kiểm tra role - chỉ manager mới có quyền
		if user.Role != usercol.RoleManager {
			code := response.ErrorResponse("Only managers can remove employees")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Kiểm tra employee có thuộc team của manager không
		belongs, err := teammembercol.CheckEmployeeBelongsToManager(c.Request.Context(), user.GetIDString(), employeeID)
		if err != nil {
			logger.Err(err).Msg("failed to check employee relationship")
			code := response.ErrorResponse("Failed to check employee relationship")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		if !belongs {
			code := response.ErrorResponse("Employee not found in your team")
			c.JSON(http.StatusNotFound, code)
			c.Abort()
			return
		}

		// Xóa team member relationship
		err = teammembercol.SoftDeleteByManagerAndEmployee(c.Request.Context(), user.GetIDString(), employeeID)
		if err != nil {
			logger.Err(err).Msg("failed to remove employee from team")
			code := response.ErrorResponse("Failed to remove employee from team")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(map[string]string{
			"message": "Employee removed from team successfully",
		}))
	}
}

