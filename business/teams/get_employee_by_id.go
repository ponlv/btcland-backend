package teams

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEmployeeByID() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][get_employee_by_id]")

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
			code := response.ErrorResponse("Only managers can view employee details")
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

		// Lấy thông tin employee
		employee, err := usercol.FindWithUserID(c.Request.Context(), employeeID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Employee not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get employee")
			code := response.ErrorResponse("Failed to get employee")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy thông tin team member
		teamMember, err := teammembercol.FindByManagerAndEmployee(c.Request.Context(), user.GetIDString(), employeeID)
		if err != nil {
			logger.Err(err).Msg("failed to get team member info")
		}

		responseData := map[string]interface{}{
			"id":             employee.GetIDString(),
			"full_name":      employee.FullName,
			"email":          employee.Email,
			"avatar":         employee.Avatar,
			"phone_number":   employee.PhoneNumber,
			"role":           employee.Role,
			"joined_at":      teamMember.JoinedAt,
			"team_member_id": teamMember.GetIDString(),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}
