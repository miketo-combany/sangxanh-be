package controller

import (
	"SangXanh/pkg/common/api"
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
	g.POST("/login", c.Login)
	g.POST("/refresh", c.Refresh)
	g.GET("/current-user", c.CurrentUser, c.auth)
}

func (c *authController) Login(e echo.Context) error {
	return api.Execute(e, c.authService.Login)
}

func (c *authController) Refresh(e echo.Context) error {
	return api.Execute(e, c.authService.Refresh)
}

func (c *authController) CurrentUser(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.authService.GetCurrentUser(ctx)
	})
}
