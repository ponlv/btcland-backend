package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func Base64ImageToReader(base64String, defaultFilename string) (io.Reader, string, string, error) {
	// Khởi tạo giá trị mặc định
	var data, mimeType string
	filename := defaultFilename

	// Tách phần header nếu chuỗi là Data URL
	if strings.Contains(base64String, ";base64,") {
		parts := strings.SplitN(base64String, ";base64,", 2)
		if len(parts) != 2 {
			return nil, "", "", fmt.Errorf("invalid base64 format")
		}
		data = parts[1]

		// Lấy MIME type từ header
		header := parts[0]
		mimeParts := strings.Split(header, ":")
		if len(mimeParts) > 1 {
			mimeType = mimeParts[1]
			// Suy ra đuôi tệp từ MIME type
			switch mimeType {
			case "image/png":
				filename = ensureExtension(defaultFilename, ".png")
			case "image/jpeg", "image/jpg":
				filename = ensureExtension(defaultFilename, ".jpg")
			case "image/gif":
				filename = ensureExtension(defaultFilename, ".gif")
			default:
				filename = ensureExtension(defaultFilename, ".jpg")
			}
		} else {
			return nil, "", "", fmt.Errorf("invalid MIME type format")
		}
	} else {
		// Nếu không có header, giả sử là chuỗi Base64 thô
		data = base64String
		mimeType = "application/octet-stream" // MIME type mặc định cho dữ liệu thô
		filename = ensureExtension(defaultFilename, ".jpg")
	}

	// Giải mã chuỗi Base64
	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Tạo io.Reader từ dữ liệu
	reader := bytes.NewBuffer(imgData)

	return reader, data, filename, nil
}

// ensureExtension đảm bảo tên tệp có đuôi phù hợp
func ensureExtension(filename, defaultExt string) string {
	if filepath.Ext(filename) == "" {
		return filename + defaultExt
	}
	return filename
}
