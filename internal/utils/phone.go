package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

// isValidPhoneNumber kiểm tra xem chuỗi có phải là số điện thoại Việt Nam hợp lệ
func IsValidPhoneNumber(phone string) bool {
	// Loại bỏ khoảng trắng, dấu gạch ngang hoặc các ký tự không cần thiết
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Kiểm tra độ dài (phải là 10 hoặc 11 chữ số)
	if len(phone) != 10 && len(phone) != 11 {
		return false
	}

	// Kiểm tra bắt đầu bằng 0
	if !strings.HasPrefix(phone, "0") {
		return false
	}

	// Kiểm tra chỉ chứa số
	matched, err := regexp.MatchString(`^\d+$`, phone)
	if err != nil || !matched {
		return false
	}

	return true
}

func IsEmailValid(email string) bool {
	// Remove leading and trailing whitespace
	email = strings.TrimSpace(email)

	// Check if email is empty
	if email == "" {
		return false
	}

	// Use net/mail to parse the email address
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Additional regex check for common email format (optional but recommended)
	// This ensures the email follows a basic pattern like local@domain.tld
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	return true
}

// convertPhoneToInternational chuyển đổi số điện thoại sang định dạng quốc tế với mã vùng +84
func ConvertPhoneToInternational(phone string) (string, error) {

	// Loại bỏ khoảng trắng, dấu gạch ngang
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Thay số 0 đầu tiên bằng +84
	result := "84" + phone[1:]
	return result, nil
}
