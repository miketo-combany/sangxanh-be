package controller

import (
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
	userService service.UserService
}

func NewUserController(di do.Injector) (api.Controller, error) {
	return &userController{
		userService: do.MustInvoke[service.UserService](di),
	}, nil
}

// --------------------------------------------------------------------
// Route registration
// --------------------------------------------------------------------

func (c *userController) Register(g *echo.Group) {
	g = g.Group("/user")

	g.GET("", c.List)
	g.POST("/register", c.Create)
	g.PUT("/update", c.Update)
	g.PUT("/address", c.Address)

	g.GET("/:id", c.GetById) // ← NEW
}

func (c *userController) GetById(e echo.Context) error {
	id := e.Param("id")

	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.userService.GetUserById(ctx, id)
	})
}

// GET /user
func (c *userController) List(e echo.Context) error {
	return api.Execute[dto.ListUser](e, func(ctx context.Context, req dto.ListUser) (api.Response, error) {
		return c.userService.ListUser(ctx, req)
	})
}

// POST /user/register
func (c *userController) Create(e echo.Context) error {
	return api.Execute(e, c.userService.Register)
}

// PUT /user/update
func (c *userController) Update(e echo.Context) error {
	return api.Execute(e, c.userService.UpdateUser)
}

// PATCH /user/address
func (c *userController) Address(e echo.Context) error {
	return api.Execute(e, c.userService.UpdateUserAddress)
}
