package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type categoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(di do.Injector) (api.Controller, error) {
	return &categoryController{
		categoryService: do.MustInvoke[service.CategoryService](di),
	}, nil
}

func (u *categoryController) Register(g *echo.Group) {
	g = g.Group("/category")
	g.GET("", u.List)
	g.POST("/create", u.Create)
	g.PUT("/update", u.Update)
}

func (u *categoryController) List(e echo.Context) error {
	name := e.QueryParam("name") // Get "name" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return u.categoryService.ListCategories(ctx, name)
	})
}

func (u *categoryController) Create(e echo.Context) error {
	return api.Execute(e, u.categoryService.CreateCategory)
}

func (u *categoryController) Update(e echo.Context) error {
	return api.Execute(e, u.categoryService.UpdateCategory)
}
