package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

var (
	ErrInvalid   = errors.New("invalid token")
	ErrExpired   = errors.New("token expired")
	ErrMalformed = errors.New("malformed token")
	ErrBadAlg    = errors.New("bad signing method")
)

func VerifyJWT(tokenString, jwtKey string) (*jwt.Token, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Only accept HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrBadAlg
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		// Detect common errors via error message string (fallback safe way)
		msg := err.Error()

		switch {
		case strings.Contains(msg, "token is expired"):
			return nil, ErrExpired
		case strings.Contains(msg, "token is malformed"):
			return nil, ErrMalformed
		case strings.Contains(msg, "signature is invalid"):
			return nil, ErrInvalid
		default:
			return nil, fmt.Errorf("jwt parse error: %w", err)
		}
	}

	if !token.Valid {
		return nil, ErrInvalid
	}

	return token, nil
}
