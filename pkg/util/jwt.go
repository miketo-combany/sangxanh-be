package util

import (
	"SangXanh/pkg/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	jwtConfig config.JWTKey
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, userRole string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtConfig.Key)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtConfig.Key)
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	var secret []byte

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
