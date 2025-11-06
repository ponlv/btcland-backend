package profile

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

func Get() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][profile][get]")

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

		// Lấy lại thông tin user từ database để đảm bảo có dữ liệu mới nhất
		currentUser, err := usercol.FindWithUserID(c.Request.Context(), user.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get user profile")
			code := response.ErrorResponse("Failed to get profile")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Trả về thông tin profile (không bao gồm password)
		profileData := map[string]interface{}{
			"id":              currentUser.GetIDString(),
			"full_name":       currentUser.FullName,
			"email":           currentUser.Email,
			"phone_number":    currentUser.PhoneNumber,
			"avatar":          currentUser.Avatar,
			"role":            currentUser.Role,
			"is_verify_phone": currentUser.IsVerifyPhone,
			"is_verify_email": currentUser.IsVerifyEmail,
			"created_at":      currentUser.CreatedAt,
			"updated_at":      currentUser.UpdatedAt,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(profileData))
	}
}

