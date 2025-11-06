package departments

import (
	"errors"
	"net/http"
	"strconv"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/teammembercol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListEmployees() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][list_employees]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
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
			code := response.ErrorResponse("Only leaders can view team employees")
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

		// Nếu team chưa có manager, trả về danh sách rỗng
		if team.ManagerID == "" {
			responseData := map[string]interface{}{
				"team_id":    team.GetIDString(),
				"team_name":  team.Name,
				"data":       []interface{}{},
				"total":      int64(0),
				"page":       1,
				"limit":      20,
				"total_page": 0,
			}
			c.JSON(http.StatusOK, response.SuccessResponse(responseData))
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

		// Lấy danh sách nhân viên trong team (thông qua TeamMember với manager_id)
		teamMembers, count, err := teammembercol.FindByManagerID(c.Request.Context(), team.ManagerID, findOptions)
		if err != nil {
			logger.Err(err).Msg("failed to list team employees")
			code := response.ErrorResponse("Failed to list team employees")
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
			"team_id":    team.GetIDString(),
			"team_name":  team.Name,
			"data":       employees,
			"total":      count,
			"page":       page,
			"limit":      limit,
			"total_page": (count + int64(limit) - 1) / int64(limit),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

