package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
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

func (c *categoryController) Register(g *echo.Group) {
	g = g.Group("/category")
	g.GET("", c.List)
	g.GET("/:id", c.GetById)
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update)
	g.DELETE("/delete", c.Delete)
}

func (c *categoryController) List(e echo.Context) error {
	// Anything that is *not* part of the DTO (here: the free‑text filter "name")
	// can be taken directly from the query‑string.
	name := e.QueryParam("name")

	// Let api.Execute bind the query params (`page`, `limit`, …) into
	// dto.ListCategory and run validation tags automatically.
	return api.Execute[dto.ListCategory](e, func(ctx context.Context, req dto.ListCategory) (api.Response, error) {
		return c.categoryService.ListCategories(ctx, req, name)
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
