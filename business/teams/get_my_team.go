package teams

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

func GetMyTeam() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][get_my_team]")

	return func(c *gin.Context) {
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
			code := response.ErrorResponse("Only managers can view team information")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Lấy tất cả nhân viên trong team
		teamMembers, count, err := teammembercol.FindByManagerID(c.Request.Context(), user.GetIDString(), nil)
		if err != nil {
			logger.Err(err).Msg("failed to get team members")
			code := response.ErrorResponse("Failed to get team information")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy thông tin chi tiết của từng employee
		employees := make([]map[string]interface{}, 0)
		for _, tm := range teamMembers {
			employee, err := usercol.FindWithUserID(c.Request.Context(), tm.EmployeeID)
			if err != nil {
				logger.Err(err).Msgf("failed to get employee: %s", tm.EmployeeID)
				continue
			}

			employees = append(employees, map[string]interface{}{
				"id":             employee.GetIDString(),
				"full_name":      employee.FullName,
				"email":          employee.Email,
				"avatar":         employee.Avatar,
				"role":           employee.Role,
				"joined_at":      tm.JoinedAt,
				"team_member_id": tm.GetIDString(),
			})
		}

		responseData := map[string]interface{}{
			"manager": map[string]interface{}{
				"id":        user.GetIDString(),
				"full_name": user.FullName,
				"email":     user.Email,
				"avatar":    user.Avatar,
			},
			"employees":       employees,
			"total_employees": count,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}
