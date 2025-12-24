package controller

import (
	"SangXanh/cmd/api/middleware"
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/log"
	"SangXanh/pkg/service"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type authController struct {
	authService service.AuthService
	auth        echo.MiddlewareFunc
}

func NewAuthController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &authController{
		authService: do.MustInvoke[service.AuthService](di),
		auth:        auth,
	}, nil
}

func (c *authController) Register(g *echo.Group) {
	g = g.Group("/auth")
	g.POST("/login", c.Login, middleware.RateLimiterMiddleware)
	g.POST("/refresh", c.Refresh, middleware.RateLimiterMiddleware)
	g.GET("/current-user", c.CurrentUser, c.auth, middleware.RateLimiterMiddleware)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with username/email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse "Successfully authenticated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Router /auth/login [post]
func (c *authController) Login(e echo.Context) error {
	ip := middleware.GetClientIP(e.Request())
	log.Info("Login from IP: %s", ip)
	return api.Execute(e, c.authService.Login)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.AuthResponse "Successfully refreshed token"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid or expired refresh token"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Router /auth/refresh [post]
func (c *authController) Refresh(e echo.Context) error {
	ip := middleware.GetClientIP(e.Request())
	log.Info("Login from IP: %s", ip)
	return api.Execute(e, c.authService.Refresh)
}

// CurrentUser godoc
// @Summary Get current user information
// @Description Retrieve the authenticated user's information
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserInfo "Current user information"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid or missing token"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Security BearerAuth
// @Router /auth/current-user [get]
func (c *authController) CurrentUser(e echo.Context) error {
	ip := middleware.GetClientIP(e.Request())
	log.Info("Login from IP: %s", ip)
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.authService.GetCurrentUser(ctx)
	})
}
