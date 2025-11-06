package teams

import (
	"net/http"
	"strconv"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListMyEmployees() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][list_my_employees]")

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
			code := response.ErrorResponse("Only managers can view employees")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Setup pagination
		skip := (page - 1) * limit
		findOptions := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(primitive.D{{Key: "joined_at", Value: -1}})

		// Tìm tất cả nhân viên của manager
		teamMembers, count, err := teammembercol.FindByManagerID(c.Request.Context(), user.GetIDString(), findOptions)
		if err != nil {
			logger.Err(err).Msg("failed to list employees")
			code := response.ErrorResponse("Failed to list employees")
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
				"id":         employee.GetIDString(),
				"full_name":  employee.FullName,
				"email":      employee.Email,
				"avatar":     employee.Avatar,
				"role":       employee.Role,
				"joined_at":  tm.JoinedAt,
				"team_member_id": tm.GetIDString(),
			})
		}

		responseData := map[string]interface{}{
			"data":       employees,
			"total":      count,
			"page":       page,
			"limit":      limit,
			"total_page": (count + int64(limit) - 1) / int64(limit),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

