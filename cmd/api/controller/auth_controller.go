package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type authController struct {
	authService service.AuthService
}

func NewAuthController(di do.Injector) (api.Controller, error) {
	return &authController{
		authService: do.MustInvoke[service.AuthService](di),
	}, nil
}

func (c *authController) Register(g *echo.Group) {
	g = g.Group("/auth")
	g.POST("/login", c.Login)
	g.POST("/refresh", c.Refresh)
}

func (c *authController) Login(e echo.Context) error {
	return api.Execute(e, c.authService.Login)
}

func (c *authController) Refresh(e echo.Context) error {
	return api.Execute(e, c.authService.Refresh)
}
