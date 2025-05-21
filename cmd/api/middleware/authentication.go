package middleware

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
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
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			// Parse custom claims and store in context if needed
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("user_id", claims["sub"])
				c.Set("user_role", claims["user_role"])
				c.Set("exp", claims["exp"])
			}

			return next(c)
		}
	}
}

func GetCurrentUser(c echo.Context) (api.Response, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User role not found in context")
	}

	return api.Success(util.CustomClaims{
		UserID:   userID,
		UserRole: userRole,
	}), nil
}
