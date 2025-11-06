package workconfirmations

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api/internal/plog"
	"api/internal/response"
	"api/schema/usercol"
	"api/schema/workconfirmationcol"
	"api/services/minio"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][update]")

	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			code := response.ErrorResponse("ID is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Parse multipart form
		err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			code := response.ErrorResponse("Failed to parse multipart form")
			c.JSON(http.StatusBadRequest, code)
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

		// Tìm đơn hiện tại
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

		// Kiểm tra quyền: chỉ người tạo mới được sửa
		if workConfirmation.CreatedBy != user.GetIDString() {
			code := response.ErrorResponse("Only creator can update")
			c.JSON(http.StatusForbidden, code)
			c.Abort()
			return
		}

		// Kiểm tra: chỉ được sửa khi chưa được xác nhận
		if workConfirmation.Status != workconfirmationcol.StatusPendingManager &&
			workConfirmation.Status != workconfirmationcol.StatusPendingLeader {
			code := response.ErrorResponse("Cannot update work confirmation that has been approved or rejected")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Lấy date và content từ form
		date := c.PostForm("date")
		content := c.PostForm("content")

		// Cập nhật thông tin
		if date != "" {
			// Validate date format
			_, err := time.Parse("2006-01-02", date)
			if err != nil {
				code := response.ErrorResponse("Invalid date format. Expected YYYY-MM-DD")
				c.JSON(http.StatusBadRequest, code)
				c.Abort()
				return
			}
			workConfirmation.Date = date
		}

		if content != "" {
			workConfirmation.Content = content
		}

		// Xử lý photos nếu có
		formFiles := c.Request.MultipartForm.File["photos"]
		if len(formFiles) > 0 {
			// Upload các file mới lên MinIO và thay thế photos cũ
			photos := make([]workconfirmationcol.Photo, 0)
			bucket := "images"
			userID := user.GetIDString()

			for _, fileHeader := range formFiles {
				// Mở file
				file, err := fileHeader.Open()
				if err != nil {
					logger.Err(err).Msgf("failed to open file: %s", fileHeader.Filename)
					continue
				}

				// Đảm bảo file luôn được đóng
				func() {
					defer file.Close()

					// Validate file type (chỉ cho phép image)
					contentType := fileHeader.Header.Get("Content-Type")
					if !strings.HasPrefix(contentType, "image/") {
						logger.Warn().Msgf("invalid file type: %s", contentType)
						return
					}

					// Tạo object key: work-confirmations/{user_id}/{timestamp}-{filename}
					timestamp := time.Now().Unix()
					filename := fileHeader.Filename
					// Sanitize filename
					filename = strings.ReplaceAll(filename, " ", "_")
					filename = strings.ReplaceAll(filename, "/", "_")
					objectKey := fmt.Sprintf("work-confirmations/%s/%d-%s", userID, timestamp, filename)

					// Upload lên MinIO
					_, err = minio.PutObject(bucket, objectKey, file)
					if err != nil {
						logger.Err(err).Msgf("failed to upload file to MinIO: %s", filename)
						return
					}

					// Tạo URL
					photoURL := fmt.Sprintf("/%s/%s", bucket, objectKey)

					photos = append(photos, workconfirmationcol.Photo{
						URL:        photoURL,
						Filename:   filename,
						UploadedAt: time.Now(),
					})
				}()
			}

			// Chỉ cập nhật photos nếu có ít nhất 1 file upload thành công
			if len(photos) > 0 {
				workConfirmation.Photos = photos
			}
		}

		// Lưu cập nhật
		_, err = workconfirmationcol.Update(c.Request.Context(), workConfirmation)
		if err != nil {
			logger.Err(err).Msg("failed to update work confirmation")
			code := response.ErrorResponse("Failed to update work confirmation")
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
