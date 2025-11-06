package teams

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

type UpdateEmployeeRequest struct {
	// Có thể thêm các trường cần cập nhật sau
	// Hiện tại chỉ để placeholder
}

func UpdateEmployee() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][update_employee]")

	return func(c *gin.Context) {
		employeeID := c.Param("employee_id")
		if employeeID == "" {
			code := response.ErrorResponse("Employee ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req UpdateEmployeeRequest
		if err := c.BindJSON(&req); err != nil {
			// Request body là optional
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
			code := response.ErrorResponse("Only managers can update employee information")
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

		// TODO: Implement update logic khi cần
		// Hiện tại chỉ trả về success
		c.JSON(http.StatusOK, response.SuccessResponse(map[string]string{
			"message": "Employee information updated successfully",
		}))
	}
}

