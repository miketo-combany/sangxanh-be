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

type categoryController struct {
	categoryService service.CategoryService
	middleware      echo.MiddlewareFunc
}

func NewCategoryController(di do.Injector, middleware echo.MiddlewareFunc) (api.Controller, error) {
	return &categoryController{
		categoryService: do.MustInvoke[service.CategoryService](di),
		middleware:      middleware,
	}, nil
}

func (c *categoryController) Register(g *echo.Group) {
	g = g.Group("/category")
	g.GET("", c.List)
	g.GET("/:id", c.GetById)
	g.POST("/create", c.Create, c.middleware, middleware.RequireRoles("admin"))
	g.PUT("/update", c.Update, c.middleware, middleware.RequireRoles("admin"))
	g.DELETE("/delete", c.Delete, c.middleware, middleware.RequireRoles("admin"))
}

func (c *categoryController) List(e echo.Context) error {
	return api.Execute[dto.ListCategory](e, func(ctx context.Context, req dto.ListCategory) (api.Response, error) {
		return c.categoryService.ListCategories(ctx, req)
	})
}

func (c *categoryController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.ListCategoryById(ctx, id)
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
