package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// RequireRoles ensures the user has one of the allowed roles.
// Assumes the AuthenticationMiddleware has already stored "user_role" in the context.
func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || role == "" {
				return echo.NewHTTPError(http.StatusForbidden, "Missing user role")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Access denied: insufficient role")
		}
	}
}
