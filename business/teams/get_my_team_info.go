package teams

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

func GetMyTeamInfo() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][teams][get_my_team_info]")

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

		// Tìm team mà manager đang quản lý
		teams, _, err := teamcol.FindByManagerID(c.Request.Context(), user.GetIDString(), nil)
		if err != nil {
			logger.Err(err).Msg("failed to get team")
			code := response.ErrorResponse("Failed to get team information")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Nếu manager chưa được gán vào team nào
		if len(teams) == 0 {
			responseData := map[string]interface{}{
				"team": nil,
				"message": "Bạn chưa được gán vào phòng ban nào",
			}
			c.JSON(http.StatusOK, response.SuccessResponse(responseData))
			return
		}

		// Lấy team đầu tiên (manager chỉ quản lý 1 team)
		team := teams[0]

		// Format response
		teamData := map[string]interface{}{
			"id":          team.GetIDString(),
			"name":        team.Name,
			"description": team.Description,
			"manager_id":  team.ManagerID,
			"created_at":  team.CreatedAt,
			"updated_at":  team.UpdatedAt,
		}

		// Lấy thông tin manager
		manager, err := usercol.FindWithUserID(c.Request.Context(), team.ManagerID)
		if err == nil && manager != nil {
			teamData["manager"] = map[string]interface{}{
				"id":        manager.GetIDString(),
				"full_name": manager.FullName,
				"email":     manager.Email,
				"avatar":    manager.Avatar,
			}
		}

		responseData := map[string]interface{}{
			"team": teamData,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

