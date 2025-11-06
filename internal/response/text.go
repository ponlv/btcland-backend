package response

import "strings"

var errorTextMap = map[string]int{
	// Authentication errors
	"TOKEN_EXPIRED": 401,
	// Validation errors
	"INVALID_PARAM": 400,
	// Account status errors
	"ACCOUNT_NOT_VERIFY_PHONE": 404,
	// Conflict errors - user already exists
	"ACCOUNT_EXIST": 409,
	"PHONE_EXIST":   409,
	// Server errors
	"SERVER_ERROR": 500,
	// Service unavailable errors
	"SERVICE_UNAVAILABLE": 503,
}

func GetCode(e string) int {

	for key, value := range errorTextMap {
		if strings.Contains(e, key) {
			return value
		}
	}
	return 400
}
