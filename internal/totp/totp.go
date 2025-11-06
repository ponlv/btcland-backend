package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"strings"
	"time"
)

// TOTP represents a TOTP generator
type TOTP struct {
	Secret string
	Digits int
	Period int
}

// NewTOTP creates a new TOTP instance
func NewTOTP(secret string) *TOTP {
	return &TOTP{
		Secret: secret,
		Digits: 6,
		Period: 30,
	}
}

// GenerateCode generates a TOTP code for the current time
func (t *TOTP) GenerateCode() (string, error) {
	return t.GenerateCodeAt(time.Now())
}

// GenerateCodeAt generates a TOTP code for a specific time
func (t *TOTP) GenerateCodeAt(timestamp time.Time) (string, error) {
	counter := uint64(timestamp.Unix()) / uint64(t.Period)
	return t.GenerateCodeForCounter(counter)
}

// GenerateCodeForCounter generates a TOTP code for a specific counter
func (t *TOTP) GenerateCodeForCounter(counter uint64) (string, error) {
	// Decode the secret
	secret, err := base32.StdEncoding.DecodeString(t.Secret)
	if err != nil {
		return "", err
	}

	// Convert counter to bytes
	counterBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		counterBytes[i] = byte(counter)
		counter >>= 8
	}

	// Calculate HMAC-SHA1
	h := hmac.New(sha1.New, secret)
	h.Write(counterBytes)
	hmacResult := h.Sum(nil)

	// Dynamic truncation
	offset := hmacResult[19] & 0xf
	code := ((int(hmacResult[offset]) & 0x7f) << 24) |
		((int(hmacResult[offset+1]) & 0xff) << 16) |
		((int(hmacResult[offset+2]) & 0xff) << 8) |
		(int(hmacResult[offset+3]) & 0xff)

	code = code % int(math.Pow10(t.Digits))

	// Format the code with leading zeros
	format := fmt.Sprintf("%%0%dd", t.Digits)
	return fmt.Sprintf(format, code), nil
}

// VerifyCode verifies a TOTP code
func (t *TOTP) VerifyCode(code string) bool {
	now := time.Now()

	// Check current window
	if t.verifyCodeAt(code, now) {
		return true
	}

	// Check previous window (for clock skew tolerance)
	if t.verifyCodeAt(code, now.Add(-time.Duration(t.Period)*time.Second)) {
		return true
	}

	// Check next window (for clock skew tolerance)
	if t.verifyCodeAt(code, now.Add(time.Duration(t.Period)*time.Second)) {
		return true
	}

	return false
}

// verifyCodeAt verifies a code at a specific time
func (t *TOTP) verifyCodeAt(code string, timestamp time.Time) bool {
	generatedCode, err := t.GenerateCodeAt(timestamp)
	if err != nil {
		return false
	}

	return generatedCode == code
}

// GenerateSecret generates a random secret key
func GenerateSecret() (string, error) {
	// This is a simplified version - in production, use crypto/rand
	secret := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	return secret, nil
}

// FormatSecret formats a secret key for display
func FormatSecret(secret string) string {
	// Remove padding
	secret = strings.TrimRight(secret, "=")

	// Add spaces every 4 characters for readability
	result := ""
	for i, char := range secret {
		if i > 0 && i%4 == 0 {
			result += " "
		}
		result += string(char)
	}

	return result
}

// GenerateQRCodeURL generates a QR code URL for TOTP setup
func GenerateQRCodeURL(issuer, accountName, secret string) string {
	// Encode the parameters
	encodedIssuer := strings.ReplaceAll(issuer, " ", "%20")
	encodedAccount := strings.ReplaceAll(accountName, " ", "%20")
	encodedSecret := strings.ReplaceAll(secret, " ", "%20")

	// Create the TOTP URL
	url := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		encodedIssuer, encodedAccount, encodedSecret, encodedIssuer)

	return url
}
