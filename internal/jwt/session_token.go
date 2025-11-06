package jwt

import (
	"api/internal/timer"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// SessionToken represents the structure of the session token
type SessionToken struct {
	SessionId string `json:"session_id"`
	Type      string `json:"type"`
	jwt.RegisteredClaims
}

const (
	SessionTypeForgotPass        string = "forgot_pass"
	SessionTypeChangePass        string = "change_pass"
	SessionTypeRegister          string = "register"
	SessionTypeChangePhone       string = "change_phone"
	SessionTypeChangeEmail       string = "change_email"
	SessionTypeChangeEmailVerify string = "change_email_verify"
	SessionTypeChangePhoneVerify string = "change_phone_verify"
)

// GenerateSessionToken generates a session token with the given parameters
func GenerateSessionToken(keySign string, sessionId, sessionType, issuer string, expired int) (string, error) {

	// Check if the key is empty
	if keySign == "" {
		return "", errors.New("key is empty")
	}

	signingKey := []byte(keySign)

	// Create the claims
	claims := SessionToken{
		sessionId,
		sessionType,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(timer.Now().Add(time.Duration(expired) * time.Second)),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	res, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return res, nil
}

func VerifySessionToken(key, tokenString string) (*SessionToken, error) {
	if key == "" {

		return nil, errors.New("_KEY_IS_EMPTY_")
	}

	if tokenString == "" {
		return nil, errors.New("_TOKEN_IS_EMPTY_")
	}

	token, err := jwt.ParseWithClaims(tokenString, &SessionToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected contract method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		var v *jwt.ValidationError
		if errors.As(err, &v) {
			if v.Errors == jwt.ValidationErrorExpired {
				return nil, errors.New("_TOKEN_EXPIRED_")
			}
		}
		return nil, errors.New("_TOKEN_EXPIRED_")
	}
	if token == nil {
		return nil, errors.New("PARSE_TOKEN_ERROR")
	}
	claim, ok := token.Claims.(*SessionToken)
	if ok && token.Valid {
		return claim, nil
	}

	return nil, err
}
