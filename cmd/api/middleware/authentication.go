package middleware

import (
	"SangXanh/pkg/util"
	"github.com/labstack/echo/v4"
	"strings"
)

func AuthenticationMiddleware(jwtKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.ErrUnauthorized
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := util.VerifyJWT(tokenString, jwtKey)
			if err != nil {
				return echo.ErrUnauthorized
			}

			// Parse custom claims and store in context if needed
			if claims, ok := token.Claims.(*util.CustomClaims); ok {
				c.Set("user_id", claims.UserID)
				c.Set("user_role", claims.UserRole)
			}

			return next(c)
		}
	}
}
