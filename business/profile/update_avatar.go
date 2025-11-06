package profile

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"
	"api/services/minio"

	"github.com/gin-gonic/gin"
)

func UpdateAvatar() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][profile][update_avatar]")

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

		// Parse multipart form
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			code := response.ErrorResponse("Failed to parse multipart form")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Lấy file avatar từ form
		formFiles := c.Request.MultipartForm.File["avatar"]
		if len(formFiles) == 0 {
			code := response.ErrorResponse("avatar file is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Chỉ lấy file đầu tiên
		fileHeader := formFiles[0]

		// Mở file
		file, err := fileHeader.Open()
		if err != nil {
			logger.Err(err).Msgf("failed to open file: %s", fileHeader.Filename)
			code := response.ErrorResponse("Failed to open file")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}
		defer file.Close()

		// Validate file type (chỉ cho phép image)
		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			code := response.ErrorResponse("Invalid file type. Only images are allowed")
			c.JSON(http.StatusBadRequest, code)
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

		// Xóa avatar cũ nếu có (tùy chọn, có thể bỏ qua nếu muốn giữ lại)
		// TODO: Có thể thêm logic xóa avatar cũ từ MinIO nếu cần

		// Tạo object key: avatars/{user_id}/{timestamp}-{filename}
		timestamp := time.Now().Unix()
		filename := fileHeader.Filename
		// Sanitize filename
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ReplaceAll(filename, "/", "_")
		objectKey := fmt.Sprintf("avatars/%s/%d-%s", user.GetIDString(), timestamp, filename)

		// Upload lên MinIO
		bucket := "images"
		_, err = minio.PutObject(bucket, objectKey, file)
		if err != nil {
			logger.Err(err).Msgf("failed to upload avatar to MinIO: %s", filename)
			code := response.ErrorResponse("Failed to upload avatar")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Tạo URL
		avatarURL := fmt.Sprintf("/%s/%s", bucket, objectKey)

		// Cập nhật avatar trong database
		currentUser.Avatar = avatarURL
		_, err = usercol.Update(c.Request.Context(), currentUser)
		if err != nil {
			logger.Err(err).Msg("failed to update avatar")
			code := response.ErrorResponse("Failed to update avatar")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại user đã cập nhật
		updated, err := usercol.FindWithUserID(c.Request.Context(), user.GetIDString())
		if err != nil {
			logger.Err(err).Msg("failed to get updated profile")
			code := response.ErrorResponse("Avatar updated but failed to retrieve")
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
