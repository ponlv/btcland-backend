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

func GetByID() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][get_by_id]")

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
			code := response.ErrorResponse("Only leaders can view team details")
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

		// Lấy thông tin manager nếu có
		teamData := map[string]interface{}{
			"id":          team.GetIDString(),
			"name":        team.Name,
			"description": team.Description,
			"manager_id":  team.ManagerID,
			"created_at":  team.CreatedAt,
			"updated_at":  team.UpdatedAt,
		}

		if team.ManagerID != "" {
			manager, err := usercol.FindWithUserID(c.Request.Context(), team.ManagerID)
			if err == nil && manager != nil {
				teamData["manager"] = map[string]interface{}{
					"id":        manager.GetIDString(),
					"full_name": manager.FullName,
					"email":     manager.Email,
					"avatar":    manager.Avatar,
				}
			}
		}

		// Lấy danh sách nhân viên trong team (thông qua TeamMember với manager_id)
		if team.ManagerID != "" {
			teamMembers, _, err := teammembercol.FindByManagerID(c.Request.Context(), team.ManagerID, nil)
			if err == nil {
				employees := make([]map[string]interface{}, 0)
				for _, tm := range teamMembers {
					employee, err := usercol.FindWithUserID(c.Request.Context(), tm.EmployeeID)
					if err == nil && employee != nil {
						employees = append(employees, map[string]interface{}{
							"id":             employee.GetIDString(),
							"full_name":      employee.FullName,
							"email":           employee.Email,
							"avatar":          employee.Avatar,
							"role":            employee.Role,
							"joined_at":       tm.JoinedAt,
							"team_member_id":  tm.GetIDString(),
						})
					}
				}
				teamData["employees"] = employees
			}
		} else {
			teamData["employees"] = []interface{}{}
		}

		c.JSON(http.StatusOK, response.SuccessResponse(teamData))
	}
}

