package controller

import (
	"SangXanh/cmd/api/middleware"
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/service"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

// --------------------------------------------------------------------
// Controller struct & factory
// --------------------------------------------------------------------

type userController struct {
	userService    service.UserService
	authMiddleware echo.MiddlewareFunc
}

func NewUserController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &userController{
		userService:    do.MustInvoke[service.UserService](di),
		authMiddleware: auth,
	}, nil
}

// --------------------------------------------------------------------
// Route registration
// --------------------------------------------------------------------

func (c *userController) Register(g *echo.Group) {
	g = g.Group("/user")

	g.GET("", c.List)
	g.POST("/register", c.Create, middleware.RateLimiterMiddleware)
	g.PUT("/update", c.Update, c.authMiddleware, middleware.RateLimiterMiddleware)
	g.PUT("/address", c.Address, c.authMiddleware, middleware.RateLimiterMiddleware)
	g.PUT("/change-password", c.ChangePassword, c.authMiddleware, middleware.RateLimiterMiddleware)
	g.PUT("/send-magic-link", c.SendMagicLink, middleware.RateLimiterMiddleware)
	g.PUT("/forgot-password", c.ForgotPassword, c.authMiddleware, middleware.RateLimiterMiddleware)
	g.PUT("/update-status", c.UpdateStatus, c.authMiddleware, middleware.RequireRoles("admin"))

	g.GET("/:id", c.GetById) // ← NEW
}

// GetById godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserInfo "User information"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /user/{id} [get]
func (c *userController) GetById(e echo.Context) error {
	id := e.Param("id")

	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.userService.GetUserById(ctx, id)
	})
}

// List godoc
// @Summary List all users
// @Description Get a paginated list of users with optional filters
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Param name query string false "Filter by name"
// @Success 200 {object} map[string]interface{} "Paginated list of users"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /user [get]
func (c *userController) List(e echo.Context) error {
	return api.Execute[dto.ListUser](e, func(ctx context.Context, req dto.ListUser) (api.Response, error) {
		return c.userService.ListUser(ctx, req)
	})
}

// Create godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} dto.UserInfo "Created user"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Router /user/register [post]
func (c *userController) Create(e echo.Context) error {
	return api.Execute(e, c.userService.Register)
}

// Update godoc
// @Summary Update user information
// @Description Update authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UserUpdateRequest true "User update data"
// @Success 200 {object} dto.UserInfo "Updated user"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Security BearerAuth
// @Router /user/update [put]
func (c *userController) Update(e echo.Context) error {
	return api.Execute(e, c.userService.UpdateUser)
}

// Address godoc
// @Summary Update user address
// @Description Update the authenticated user's addresses
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UserUpdateAddressRequest true "Address data"
// @Success 200 {object} dto.UserInfo "Updated user with new addresses"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Security BearerAuth
// @Router /user/address [put]
func (c *userController) Address(e echo.Context) error {
	return api.Execute(e, c.userService.UpdateUserAddress)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.ChangePassword true "Password change data"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized or invalid old password"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Security BearerAuth
// @Router /user/change-password [put]
func (c *userController) ChangePassword(e echo.Context) error {
	return api.Execute(e, c.userService.ChangePassword)
}

// SendMagicLink godoc
// @Summary Send magic link for passwordless login
// @Description Send a magic link to user's email for passwordless authentication
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Email address"
// @Success 200 {object} map[string]interface{} "Magic link sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Router /user/send-magic-link [put]
func (c *userController) SendMagicLink(e echo.Context) error {
	return api.Execute(e, c.userService.SendMagicLink)
}

// ForgotPassword godoc
// @Summary Reset user password
// @Description Send password reset instructions to user's email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Email address"
// @Success 200 {object} map[string]interface{} "Password reset email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 429 {object} map[string]interface{} "Too many requests"
// @Security BearerAuth
// @Router /user/forgot-password [put]
func (c *userController) ForgotPassword(e echo.Context) error {
	return api.Execute(e, c.userService.ForgotPassword)
}

// UpdateStatus godoc
// @Summary Update user status
// @Description Update a user's status (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.UserUpdateData true "User status data"
// @Success 200 {object} dto.UserInfo "Updated user"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Security BearerAuth
// @Router /user/update-status [put]
func (c *userController) UpdateStatus(e echo.Context) error {
	return api.Execute(e, c.userService.UpdateStatus)
}
