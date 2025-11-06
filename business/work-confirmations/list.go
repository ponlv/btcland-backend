package workconfirmations

import (
	"net/http"
	"strconv"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teammembercol"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func List() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][list]")

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

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		status := c.Query("status")
		createdBy := c.Query("created_by")

		// Tạo filter
		filter := primitive.D{}

		// Filter theo role
		userID := user.GetIDString()
		userRole := user.Role

		if createdBy != "" {
			// Nếu có created_by trong query, filter theo đó
			// Nhưng vẫn cần kiểm tra quyền truy cập
			if userRole == usercol.RoleEmployee {
				// Employee chỉ có thể xem đơn của chính mình
				if createdBy != userID {
					code := response.ErrorResponse("Access denied")
					c.JSON(http.StatusForbidden, code)
					c.Abort()
					return
				}
			} else if userRole == usercol.RoleManager {
				// Manager chỉ có thể xem đơn của nhân viên trong team hoặc của chính mình
				if createdBy != userID {
					belongs, err := teammembercol.CheckEmployeeBelongsToManager(c.Request.Context(), userID, createdBy)
					if err != nil || !belongs {
						code := response.ErrorResponse("Access denied")
						c.JSON(http.StatusForbidden, code)
						c.Abort()
						return
					}
				}
			}
			// Leader và Assistant Director có thể xem tất cả
			filter = append(filter, primitive.E{Key: "created_by", Value: createdBy})
		} else {
			// Không có created_by, filter theo role
			if userRole == usercol.RoleEmployee {
				// Employee: chỉ thấy đơn của mình
				filter = append(filter, primitive.E{Key: "created_by", Value: userID})
			} else if userRole == usercol.RoleManager {
				// Manager: thấy đơn của nhân viên trong team và của chính mình
				// Lấy danh sách employee IDs trong team
				teamMembers, _, err := teammembercol.FindByManagerID(c.Request.Context(), userID, nil)
				if err != nil {
					logger.Err(err).Msg("failed to get team members")
					// Nếu lỗi, chỉ lấy đơn của chính manager
					filter = append(filter, primitive.E{Key: "created_by", Value: userID})
				} else {
					// Tạo danh sách user IDs (bao gồm cả manager)
					userIDs := []interface{}{userID}
					for _, tm := range teamMembers {
						userIDs = append(userIDs, tm.EmployeeID)
					}
					filter = append(filter, primitive.E{Key: "created_by", Value: primitive.D{{Key: "$in", Value: userIDs}}})
				}
			} else if userRole == usercol.RoleLeader || userRole == usercol.RoleAssistantDirector {
				// Leader và Assistant Director: thấy tất cả đơn (không filter created_by)
			} else {
				// Role không xác định, mặc định chỉ thấy đơn của mình
				filter = append(filter, primitive.E{Key: "created_by", Value: userID})
			}
		}

		// Filter theo status
		if status != "" {
			filter = append(filter, primitive.E{Key: "status", Value: workconfirmationcol.WorkConfirmationStatus(status)})
		}

		// Setup pagination
		skip := (page - 1) * limit
		findOptions := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(primitive.D{{Key: "created_at", Value: -1}})

		// Query
		results, count, err := workconfirmationcol.FindWithFilter(c.Request.Context(), filter, findOptions)
		if err != nil {
			logger.Err(err).Msg("failed to list work confirmations")
			code := response.ErrorResponse("Failed to list work confirmations")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		responseData := map[string]interface{}{
			"data":       results,
			"total":      count,
			"page":       page,
			"limit":      limit,
			"total_page": (count + int64(limit) - 1) / int64(limit),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}
