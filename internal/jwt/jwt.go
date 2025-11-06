package jwt

import (
	"api/internal/timer"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserId   string `json:"user_id"`
	Metadata string `json:"metadata"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(key_sign string, UserId string, metadata, issuer string, expired int) (string, error) {
	signingKey := []byte(key_sign)

	// Create the claims
	claims := CustomClaims{
		UserId,
		metadata,
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

func TokenExpiredTime(key, token_string string) float64 {
	var claims CustomClaims
	_, err := jwt.ParseWithClaims(token_string, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected contract method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err == nil {
		return 0
	}
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		//tm := time.Unix(claims.ExpiresAt, 0)
		return timer.Now().Sub(claims.ExpiresAt.Time).Seconds()
	}
	return 0
}

func VerifyJWTToken(key, tokenString string) (*CustomClaims, error) {
	if key == "" {

		return nil, errors.New("_KEY_IS_EMPTY_")
	}

	if tokenString == "" {
		return nil, errors.New("_TOKEN_IS_EMPTY_")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected contract method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors == jwt.ValidationErrorExpired {
				return nil, errors.New("_TOKEN_EXPIRED_")
			}
		}
		return nil, errors.New("_TOKEN_EXPIRED_")
	}
	if token == nil {
		return nil, errors.New("PARSE_TOKEN_ERROR")
	}

	claim, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("CLAIM_INVALID")
	}

	return claim, nil
}
