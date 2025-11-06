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

type ApproveRequest struct {
	Comment string `json:"comment"`
}

func Approve() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][approve]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		var req ApproveRequest
		if err := c.BindJSON(&req); err != nil {
			// Comment là optional, không bắt buộc
			req.Comment = ""
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

		// Xác định role của user và kiểm tra quyền
		userID := user.GetIDString()
		creatorID := workConfirmation.CreatedBy

		// Kiểm tra trạng thái và xác nhận
		if workConfirmation.Status == workconfirmationcol.StatusPendingManager {
			// Cần quản lý xác nhận
			// Kiểm tra user có phải là manager không
			if user.Role != usercol.RoleManager {
				code := response.ErrorResponse("Only managers can approve work confirmations at this stage")
				c.JSON(http.StatusForbidden, code)
				c.Abort()
				return
			}

			// Kiểm tra user có phải là manager của người tạo đơn không
			// Nếu đơn từ employee, cần kiểm tra team relationship
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
					code := response.ErrorResponse("You can only approve work confirmations from your team members")
					c.JSON(http.StatusForbidden, code)
					c.Abort()
					return
				}
			}

			err = workconfirmationcol.ApproveByManager(c.Request.Context(), id, userID, req.Comment)
			if err != nil {
				logger.Err(err).Msg("failed to approve by manager")
				code := response.ErrorResponse("Failed to approve work confirmation")
				c.JSON(http.StatusInternalServerError, code)
				c.Abort()
				return
			}
		} else if workConfirmation.Status == workconfirmationcol.StatusPendingLeader {
			// Cần lãnh đạo xác nhận
			// Kiểm tra user có phải là leader không
			if user.Role != usercol.RoleLeader {
				code := response.ErrorResponse("Only leaders can approve work confirmations at this stage")
				c.JSON(http.StatusForbidden, code)
				c.Abort()
				return
			}

			err = workconfirmationcol.ApproveByLeader(c.Request.Context(), id, userID, req.Comment)
			if err != nil {
				logger.Err(err).Msg("failed to approve by leader")
				code := response.ErrorResponse("Failed to approve work confirmation")
				c.JSON(http.StatusInternalServerError, code)
				c.Abort()
				return
			}
		} else {
			code := response.ErrorResponse("Work confirmation is not in a state that can be approved")
			c.JSON(http.StatusBadRequest, code)
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
