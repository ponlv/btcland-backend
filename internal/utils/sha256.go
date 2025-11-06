package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Sum(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

// Hàm tính SHA256 hash
func HashSHA256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Hàm tính HMAC-SHA256
func HmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Tạo contract key theo AWS specification
func GetSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kSecret := []byte("AWS4" + key)
	kDate := HmacSHA256(kSecret, dateStamp)
	kRegion := HmacSHA256(kDate, regionName)
	kService := HmacSHA256(kRegion, serviceName)
	kSigning := HmacSHA256(kService, "aws4_request")
	return kSigning
}
