package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type userController struct {
	userService service.UserService
}

func NewUserController(di do.Injector) (api.Controller, error) {
	return &userController{
		userService: do.MustInvoke[service.UserService](di),
	}, nil
}

func (u *userController) Register(g *echo.Group) {
	g = g.Group("/user")
	g.GET("", u.List)
	g.POST("", u.Create)
}

func (u *userController) List(e echo.Context) error {
	return api.Execute(e, u.userService.ListUser)
}
func (u *userController) Create(e echo.Context) error {
	return api.Execute(e, u.userService.CreateUser)
}
