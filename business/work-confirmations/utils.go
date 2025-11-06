package workconfirmations

import "strings"

// extractObjectNameFromURL trích xuất object name từ URL
func extractObjectNameFromURL(url string) string {
	// Giả sử URL có dạng: http://minio-endpoint/bucket/object-name
	// Hoặc: /bucket/object-name
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}
	// Lấy phần cuối cùng là object name
	return parts[len(parts)-1]
}

