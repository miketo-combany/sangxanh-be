package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type categoryController struct {
	categoryService service.ProductService
}

func NewCategoryController(di do.Injector) (api.Controller, error) {
	return &categoryController{
		categoryService: do.MustInvoke[service.ProductService](di),
	}, nil
}

func (c *categoryController) Register(g *echo.Group) {
	g = g.Group("/category")
	g.GET("", c.List)
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update)
	g.DELETE("/delete", c.Delete)
}

func (c *categoryController) List(e echo.Context) error {
	name := e.QueryParam("name") // Get "name" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.ListCategories(ctx, name)
	})
}

func (c *categoryController) Create(e echo.Context) error {
	return api.Execute(e, c.categoryService.CreateCategory)
}

func (c *categoryController) Update(e echo.Context) error {
	return api.Execute(e, c.categoryService.UpdateCategory)
}

func (c *categoryController) Delete(e echo.Context) error {
	id := e.QueryParam("categoryId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.DeleteCategory(ctx, id)
	})
}
