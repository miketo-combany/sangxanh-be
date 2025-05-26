package middleware

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/util"
	"context"
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

			// Extract claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := c.Request().Context()
				ctx = context.WithValue(ctx, "user_id", claims["sub"])
				ctx = context.WithValue(ctx, "user_role", claims["user_role"])

				// Replace request context
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}

func GetCurrentUser(c echo.Context) (api.Response, error) {
	ctx := c.Request().Context()

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	userRole, ok := ctx.Value("user_role").(string)
	if !ok || userRole == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User role not found in context")
	}

	return api.Success(util.CustomClaims{
		UserID:   userID,
		UserRole: userRole,
	}), nil
}
