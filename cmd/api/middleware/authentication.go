package middleware

import (
	"SangXanh/pkg/util"
	"github.com/labstack/echo/v4"
	"strings"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := util.VerifyJWT(tokenString)
		if err != nil {
			return err
		}
		return next(c)
	}
}
