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

func Create() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][work-confirmations][create]")

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
		err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			code := response.ErrorResponse("Failed to parse multipart form")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Lấy date, time và content từ form
		date := c.PostForm("date")
		startTime := c.PostForm("start_time")
		endTime := c.PostForm("end_time")
		content := c.PostForm("content")

		if date == "" {
			code := response.ErrorResponse("date is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		if startTime == "" {
			code := response.ErrorResponse("start_time is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		if endTime == "" {
			code := response.ErrorResponse("end_time is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		if content == "" {
			code := response.ErrorResponse("content is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Validate date format
		_, err = time.Parse("2006-01-02", date)
		if err != nil {
			code := response.ErrorResponse("Invalid date format. Expected YYYY-MM-DD")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Validate time format (HH:MM)
		_, err = time.Parse("15:04", startTime)
		if err != nil {
			code := response.ErrorResponse("Invalid start_time format. Expected HH:MM (24-hour format)")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		_, err = time.Parse("15:04", endTime)
		if err != nil {
			code := response.ErrorResponse("Invalid end_time format. Expected HH:MM (24-hour format)")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Validate end_time > start_time
		startTimeObj, _ := time.Parse("15:04", startTime)
		endTimeObj, _ := time.Parse("15:04", endTime)
		if !endTimeObj.After(startTimeObj) {
			code := response.ErrorResponse("end_time must be after start_time")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Xác định role của người tạo từ user
		creatorRole := user.Role
		// Nếu user chưa có role, mặc định là employee
		if creatorRole == "" {
			creatorRole = usercol.RoleEmployee
		}

		// Lấy các file từ form
		formFiles := c.Request.MultipartForm.File["photos"]
		if len(formFiles) == 0 {
			code := response.ErrorResponse("At least one photo is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Upload các file lên MinIO và tạo danh sách photos
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

				// Tạo URL (giả sử có endpoint từ env, hoặc có thể tạo sau)
				// TODO: Tạo URL đầy đủ từ MinIO endpoint
				photoURL := fmt.Sprintf("/%s/%s", bucket, objectKey)

				photos = append(photos, workconfirmationcol.Photo{
					URL:        photoURL,
					Filename:   filename,
					UploadedAt: time.Now(),
				})
			}()
		}

		if len(photos) == 0 {
			code := response.ErrorResponse("Failed to upload photos")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Xác định status ban đầu
		status := workconfirmationcol.StatusPendingManager
		if creatorRole == usercol.RoleManager {
			status = workconfirmationcol.StatusPendingLeader
		}

		// Tạo đơn mới
		workConfirmation := &workconfirmationcol.WorkConfirmation{
			CreatedBy:   user.GetIDString(),
			CreatorRole: creatorRole,
			Date:        date,
			StartTime:   startTime,
			EndTime:     endTime,
			Content:     content,
			Photos:      photos,
			Status:      status,
		}

		_, err = workconfirmationcol.Create(c.Request.Context(), workConfirmation)
		if err != nil {
			logger.Err(err).Msg("failed to create work confirmation")
			code := response.ErrorResponse("Failed to create work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		// Lấy lại đơn vừa tạo để trả về
		created, err := workconfirmationcol.FindByID(c.Request.Context(), workConfirmation.GetIDString())
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				code := response.ErrorResponse("Work confirmation not found")
				c.JSON(http.StatusNotFound, code)
				c.Abort()
				return
			}
			logger.Err(err).Msg("failed to find created work confirmation")
			code := response.ErrorResponse("Failed to retrieve work confirmation")
			c.JSON(http.StatusInternalServerError, code)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response.SuccessResponse(created))
	}
}
