package profile

import (
	"net/http"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func Update() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][profile][update]")

	return func(c *gin.Context) {
		var req UpdateRequest
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

		// Lấy lại thông tin user từ database
		currentUser, err := usercol.FindWithUserID(c.Request.Context(), user.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get user")
			code := response.ErrorResponse("Failed to get user")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Cập nhật thông tin
		if req.FullName != "" {
			currentUser.FullName = req.FullName
		}
		if req.PhoneNumber != "" {
			currentUser.PhoneNumber = req.PhoneNumber
		}

		// Lưu cập nhật
		_, err = usercol.Update(c.Request.Context(), currentUser)
		if err != nil {
			logger.Err(err).Msg("failed to update profile")
			code := response.ErrorResponse("Failed to update profile")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại user đã cập nhật
		updated, err := usercol.FindWithUserID(c.Request.Context(), user.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get updated profile")
			code := response.ErrorResponse("Profile updated but failed to retrieve")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Trả về thông tin profile (không bao gồm password)
		profileData := map[string]interface{}{
			"id":              updated.GetIDString(),
			"full_name":       updated.FullName,
			"email":           updated.Email,
			"phone_number":    updated.PhoneNumber,
			"avatar":          updated.Avatar,
			"role":            updated.Role,
			"is_verify_phone": updated.IsVerifyPhone,
			"is_verify_email": updated.IsVerifyEmail,
			"created_at":      updated.CreatedAt,
			"updated_at":      updated.UpdatedAt,
		}

		c.JSON(http.StatusOK, response.SuccessResponse(profileData))
	}
}

