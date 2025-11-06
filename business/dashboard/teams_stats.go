package dashboard

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/teammembercol"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TeamsStats() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][dashboard][teams_stats]")

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

		// Kiểm tra role - chỉ leader mới có quyền
		if user.Role != usercol.RoleLeader {
			code := response.ErrorResponse("Only leaders can view teams stats")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// Lấy tất cả teams
		teamsFilter := primitive.D{{Key: "is_delete", Value: false}}
		teams, _, err := teamcol.FindWithFilter(ctx, teamsFilter, nil)
		if err != nil {
			logger.Err(err).Msg("failed to get teams")
			code := response.ErrorResponse("Failed to get teams")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		workConfirmationsColl := workconfirmationcol.Collection()

		// Thống kê cho từng team
		teamsStats := make([]map[string]interface{}, 0)
		for _, team := range teams {
			teamStat := map[string]interface{}{
				"team_id":      team.GetIDString(),
				"team_name":    team.Name,
				"manager_id":   team.ManagerID,
				"total_employees": int64(0),
				"total_work_confirmations": int64(0),
			}

			// Lấy thông tin manager nếu có
			if team.ManagerID != "" {
				manager, err := usercol.FindWithUserID(ctx, team.ManagerID)
				if err == nil && manager != nil {
					teamStat["manager"] = map[string]interface{}{
						"id":        manager.GetIDString(),
						"full_name": manager.FullName,
						"email":     manager.Email,
					}
				}

				// Đếm số nhân viên trong team
				teamMembers, count, err := teammembercol.FindByManagerID(ctx, team.ManagerID, nil)
				if err == nil {
					teamStat["total_employees"] = count

					// Đếm số đơn xác nhận công tác của team (bao gồm manager và employees)
					userIDs := []interface{}{team.ManagerID}
					for _, tm := range teamMembers {
						userIDs = append(userIDs, tm.EmployeeID)
					}

					workConfirmationsFilter := primitive.D{
						{Key: "created_by", Value: primitive.D{{Key: "$in", Value: userIDs}}},
						{Key: "is_delete", Value: false},
					}
					workConfirmationsCount, err := workConfirmationsColl.CountWithCtx(ctx, workConfirmationsFilter)
					if err == nil {
						teamStat["total_work_confirmations"] = workConfirmationsCount
					}
				}
			}

			teamsStats = append(teamsStats, teamStat)
		}

		responseData := map[string]interface{}{
			"total_teams": len(teamsStats),
			"teams":       teamsStats,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

