package images

import (
	"net/http"
	"strings"

	"api/internal/plog"
	"api/internal/response"
	"api/services/minio"

	"github.com/gin-gonic/gin"
)

func Get() gin.HandlerFunc {
	logger := plog.NewBizLogger("[business][images][get]")

	return func(c *gin.Context) {
		// Lấy path từ URL, ví dụ: /images/images/work-confirmations/...
		// path sẽ là: images/work-confirmations/user_id/timestamp-filename
		path := c.Param("path")
		if path == "" {
			code := response.ErrorResponse("Path is required")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// Remove leading slash if exists
		path = strings.TrimPrefix(path, "/")

		// Parse path: images/work-confirmations/...
		// parts[0] should be "images" (bucket name)
		// parts[1] should be the object key (work-confirmations/user_id/timestamp-filename)
		parts := strings.SplitN(path, "/", 2)
		if len(parts) < 2 {
			code := response.ErrorResponse("Invalid path format")
			c.JSON(http.StatusBadRequest, code)
			c.Abort()
			return
		}

		// parts[0] should be "images" (bucket name)
		// parts[1] should be the object key
		bucket := parts[0]
		objectKey := parts[1]

		// Download file from MinIO
		fileData, err := minio.DownloadFile(bucket, objectKey)
		if err != nil {
			logger.Err(err).Msgf("failed to download file from MinIO: bucket=%s, key=%s", bucket, objectKey)
			code := response.ErrorResponse("Image not found")
			c.JSON(http.StatusNotFound, code)
			c.Abort()
			return
		}

		// Set content type based on file extension
		contentType := "image/jpeg" // default
		if strings.HasSuffix(strings.ToLower(objectKey), ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(strings.ToLower(objectKey), ".gif") {
			contentType = "image/gif"
		} else if strings.HasSuffix(strings.ToLower(objectKey), ".webp") {
			contentType = "image/webp"
		}

		// Set headers and return file
		c.Header("Content-Type", contentType)
		c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		c.Data(http.StatusOK, contentType, fileData)
	}
}

