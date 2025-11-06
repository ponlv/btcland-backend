package dashboard

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Stats() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][dashboard][stats]")

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
			code := response.ErrorResponse("Only leaders can view dashboard stats")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// Đếm tổng số phòng ban/team
		teamsFilter := primitive.D{{Key: "is_delete", Value: false}}
		teams, _, err := teamcol.FindWithFilter(ctx, teamsFilter, nil)
		if err != nil {
			logger.Err(err).Msg("failed to count teams")
		}
		totalTeams := int64(len(teams))

		// Đếm tổng số nhân viên
		employeesFilter := primitive.D{
			{Key: "role", Value: usercol.RoleEmployee},
			{Key: "is_delete", Value: false},
		}
		employeesColl := usercol.Collection()
		totalEmployees, err := employeesColl.CountWithCtx(ctx, employeesFilter)
		if err != nil {
			logger.Err(err).Msg("failed to count employees")
			totalEmployees = 0
		}

		// Đếm tổng số quản lý
		managersFilter := primitive.D{
			{Key: "role", Value: usercol.RoleManager},
			{Key: "is_delete", Value: false},
		}
		totalManagers, err := employeesColl.CountWithCtx(ctx, managersFilter)
		if err != nil {
			logger.Err(err).Msg("failed to count managers")
			totalManagers = 0
		}

		// Đếm tổng số đơn xác nhận công tác
		workConfirmationsFilter := primitive.D{{Key: "is_delete", Value: false}}
		workConfirmationsColl := workconfirmationcol.Collection()
		totalWorkConfirmations, err := workConfirmationsColl.CountWithCtx(ctx, workConfirmationsFilter)
		if err != nil {
			logger.Err(err).Msg("failed to count work confirmations")
			totalWorkConfirmations = 0
		}

		// Đếm đơn theo trạng thái
		pendingManagerFilter := primitive.D{
			{Key: "status", Value: workconfirmationcol.StatusPendingManager},
			{Key: "is_delete", Value: false},
		}
		pendingManagerCount, _ := workConfirmationsColl.CountWithCtx(ctx, pendingManagerFilter)

		pendingLeaderFilter := primitive.D{
			{Key: "status", Value: workconfirmationcol.StatusPendingLeader},
			{Key: "is_delete", Value: false},
		}
		pendingLeaderCount, _ := workConfirmationsColl.CountWithCtx(ctx, pendingLeaderFilter)

		approvedFilter := primitive.D{
			{Key: "status", Value: workconfirmationcol.StatusApproved},
			{Key: "is_delete", Value: false},
		}
		approvedCount, _ := workConfirmationsColl.CountWithCtx(ctx, approvedFilter)

		rejectedFilter := primitive.D{
			{Key: "status", Value: workconfirmationcol.StatusRejected},
			{Key: "is_delete", Value: false},
		}
		rejectedCount, _ := workConfirmationsColl.CountWithCtx(ctx, rejectedFilter)

		responseData := map[string]interface{}{
			"teams": map[string]interface{}{
				"total": totalTeams,
			},
			"users": map[string]interface{}{
				"total_employees": totalEmployees,
				"total_managers":  totalManagers,
			},
			"work_confirmations": map[string]interface{}{
				"total":            totalWorkConfirmations,
				"pending_manager":  pendingManagerCount,
				"pending_leader":   pendingLeaderCount,
				"approved":         approvedCount,
				"rejected":         rejectedCount,
			},
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

