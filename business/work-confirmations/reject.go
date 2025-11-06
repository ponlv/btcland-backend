package workconfirmations

import (
	"errors"
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type RejectRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func Reject() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][reject]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req RejectRequest
		if err := c.BindJSON(&req); err != nil {
			code := response.ErrorResponse(err.Error())
			c.JSON(code.Code, code)
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

		// Tìm đơn
		workConfirmation, err := workconfirmationcol.FindByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Work confirmation not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to get work confirmation")
			code := response.ErrorResponse("Failed to get work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Kiểm tra: chỉ có thể từ chối khi đơn đang chờ xác nhận
		if workConfirmation.Status != workconfirmationcol.StatusPendingManager &&
			workConfirmation.Status != workconfirmationcol.StatusPendingLeader {
			code := response.ErrorResponse("Work confirmation is not in a state that can be rejected")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Kiểm tra quyền
		userID := user.GetIDString()
		creatorID := workConfirmation.CreatedBy

		if workConfirmation.Status == workconfirmationcol.StatusPendingManager {
			// Manager chỉ có thể từ chối đơn từ nhân viên trong team
			if user.Role != usercol.RoleManager {
				code := response.ErrorResponse("Only managers can reject work confirmations at this stage")
				c.JSON(http.StatusForbidden, code)
				c.Abort()
				return
			}

			// Kiểm tra employee có thuộc team của manager không
			if workConfirmation.CreatorRole == usercol.RoleEmployee {
				belongs, err := teammembercol.CheckEmployeeBelongsToManager(c.Request.Context(), userID, creatorID)
				if err != nil {
					logger.Err(err).Msg("failed to check employee relationship")
					code := response.ErrorResponse("Failed to verify employee relationship")
					c.JSON(http.StatusInternalServerError, code)
					c.Abort()
					return
				}

				if !belongs {
					code := response.ErrorResponse("You can only reject work confirmations from your team members")
					c.JSON(http.StatusForbidden, code)
					c.Abort()
					return
				}
			}
		} else if workConfirmation.Status == workconfirmationcol.StatusPendingLeader {
			// Leader có thể từ chối đơn
			if user.Role != usercol.RoleLeader {
				code := response.ErrorResponse("Only leaders can reject work confirmations at this stage")
				c.JSON(http.StatusForbidden, code)
				c.Abort()
				return
			}
		}

		// Từ chối đơn
		err = workconfirmationcol.Reject(c.Request.Context(), id, user.GetIDString(), req.Reason)
		if err != nil {
			logger.Err(err).Msg("failed to reject work confirmation")
			code := response.ErrorResponse("Failed to reject work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại đơn đã cập nhật
		updated, err := workconfirmationcol.FindByID(c.Request.Context(), id)
		if err != nil {
			logger.Err(err).Msg("failed to get updated work confirmation")
			code := response.ErrorResponse("Failed to retrieve updated work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(updated))
	}
}
