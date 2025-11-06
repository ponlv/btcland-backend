package departments

import (
	"net/http"
	"strconv"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListUsers() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][list_users]")

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
			code := response.ErrorResponse("Only leaders can view all users")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100")) // Default 100 để lấy nhiều users
		roleFilter := c.Query("role")                              // Optional: filter by role (manager, employee)

		// Setup pagination
		skip := (page - 1) * limit
		findOptions := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(primitive.D{{Key: "created_at", Value: -1}})

		// Build filter
		filter := primitive.D{}
		if roleFilter != "" {
			// Validate role
			var role usercol.Role
			switch roleFilter {
			case "manager":
				role = usercol.RoleManager
			case "employee":
				role = usercol.RoleEmployee
			case "leader":
				role = usercol.RoleLeader
			case "assistant_director":
				role = usercol.RoleAssistantDirector
			default:
				code := response.ErrorResponse("Invalid role filter")
				c.JSON(http.StatusBadRequest, code)
				c.Abort()
				return
			}
			filter = append(filter, primitive.E{Key: "role", Value: role})
		}

		// Tìm tất cả users
		users, count, err := usercol.FindWithFilter(c.Request.Context(), filter, findOptions)
		if err != nil {
			logger.Err(err).Msg("failed to list users")
			code := response.ErrorResponse("Failed to list users")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Format response
		usersData := make([]map[string]interface{}, 0)
		for _, u := range users {
			usersData = append(usersData, map[string]interface{}{
				"id":        u.GetIDString(),
				"full_name": u.FullName,
				"email":     u.Email,
				"avatar":    u.Avatar,
				"role":      u.Role,
			})
		}

		responseData := map[string]interface{}{
			"data":       usersData,
			"total":      count,
			"page":       page,
			"limit":      limit,
			"total_page": (count + int64(limit) - 1) / int64(limit),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}
