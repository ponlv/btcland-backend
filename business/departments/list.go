package departments

import (
	"net/http"
	"strconv"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/teamcol"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func List() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][departments][list]")

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
			code := response.ErrorResponse("Only leaders can view all teams")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
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
			SetSort(primitive.D{{Key: "created_at", Value: -1}})

		// Tìm tất cả teams
		filter := primitive.D{}
		teams, count, err := teamcol.FindWithFilter(c.Request.Context(), filter, findOptions)
		if err != nil {
			logger.Err(err).Msg("failed to list teams")
			code := response.ErrorResponse("Failed to list teams")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy thông tin chi tiết cho từng team
		teamsData := make([]map[string]interface{}, 0)
		for _, team := range teams {
			teamData := map[string]interface{}{
				"id":          team.GetIDString(),
				"name":        team.Name,
				"description": team.Description,
				"manager_id":  team.ManagerID,
				"created_at":  team.CreatedAt,
				"updated_at":  team.UpdatedAt,
			}

			// Lấy thông tin manager nếu có
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

			teamsData = append(teamsData, teamData)
		}

		responseData := map[string]interface{}{
			"data":       teamsData,
			"total":      count,
			"page":       page,
			"limit":      limit,
			"total_page": (count + int64(limit) - 1) / int64(limit),
		}

		c.JSON(http.StatusOK, response.SuccessResponse(responseData))
	}
}

