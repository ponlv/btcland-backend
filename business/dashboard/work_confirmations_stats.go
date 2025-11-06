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

func WorkConfirmationsStats() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][dashboard][work_confirmations_stats]")

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
			code := response.ErrorResponse("Only leaders can view work confirmations stats")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		workConfirmationsColl := workconfirmationcol.Collection()

		// Lấy query parameter team_id nếu có
		teamID := c.Query("team_id")

		// Tạo filter cơ bản
		baseFilter := primitive.D{{Key: "is_delete", Value: false}}

		// Nếu có team_id, filter theo team (thông qua manager_id của team)
		if teamID != "" {
			team, err := teamcol.FindByID(ctx, teamID)
			if err != nil {
				logger.Err(err).Msg("failed to get team")
				code := response.ErrorResponse("Team not found")
				c.JSON(http.StatusBadRequest, code)
				c.Abort()
				return
			}

			if team.ManagerID != "" {
				// Lấy danh sách employee IDs trong team
				teamMembers, _, err := teammembercol.FindByManagerID(ctx, team.ManagerID, nil)
				if err == nil {
					userIDs := []interface{}{team.ManagerID} // Bao gồm cả manager
					for _, tm := range teamMembers {
						userIDs = append(userIDs, tm.EmployeeID)
					}
					baseFilter = append(baseFilter, primitive.E{Key: "created_by", Value: primitive.D{{Key: "$in", Value: userIDs}}})
				}
			} else {
				// Team chưa có manager, không có đơn nào
				responseData := map[string]interface{}{
					"team_id": teamID,
					"team_name": team.Name,
					"total": int64(0),
					"by_status": map[string]int64{
						"pending_manager": 0,
						"pending_leader":  0,
						"approved":        0,
						"rejected":        0,
					},
				}
				c.JSON(http.StatusOK, response.SuccessResponse(responseData))
				return
			}
		}

		// Đếm tổng số đơn
		total, err := workConfirmationsColl.CountWithCtx(ctx, baseFilter)
		if err != nil {
			logger.Err(err).Msg("failed to count work confirmations")
			total = 0
		}

		// Đếm theo trạng thái
		pendingManagerFilter := append(baseFilter, primitive.E{Key: "status", Value: workconfirmationcol.StatusPendingManager})
		pendingManagerCount, _ := workConfirmationsColl.CountWithCtx(ctx, pendingManagerFilter)

		pendingLeaderFilter := append(baseFilter, primitive.E{Key: "status", Value: workconfirmationcol.StatusPendingLeader})
		pendingLeaderCount, _ := workConfirmationsColl.CountWithCtx(ctx, pendingLeaderFilter)

		approvedFilter := append(baseFilter, primitive.E{Key: "status", Value: workconfirmationcol.StatusApproved})
		approvedCount, _ := workConfirmationsColl.CountWithCtx(ctx, approvedFilter)

		rejectedFilter := append(baseFilter, primitive.E{Key: "status", Value: workconfirmationcol.StatusRejected})
		rejectedCount, _ := workConfirmationsColl.CountWithCtx(ctx, rejectedFilter)

		responseData := map[string]interface{}{
			"total": total,
			"by_status": map[string]int64{
				"pending_manager": pendingManagerCount,
				"pending_leader":  pendingLeaderCount,
				"approved":        approvedCount,
				"rejected":        rejectedCount,
			},
		}

		// Thêm thông tin team nếu có filter theo team
		if teamID != "" {
			team, _ := teamcol.FindByID(ctx, teamID)
			if team != nil {
				responseData["team_id"] = teamID
				responseData["team_name"] = team.Name
			}
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

