package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type productController struct {
	productService service.ProductService
}

func NewProductController(di do.Injector) (api.Controller, error) {
	return &productController{
		productService: do.MustInvoke[service.ProductService](di),
	}, nil
}

func (c *productController) Register(g *echo.Group) {
	g = g.Group("/product")
	g.GET("", c.List)
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update)
	g.DELETE("/delete", c.Delete)
}

func (c *productController) List(e echo.Context) error {
	return api.Execute(e, c.productService.ListProducts)
}

func (c *productController) Create(e echo.Context) error {
	return api.Execute(e, c.productService.CreateProduct)
}

func (c *productController) Update(e echo.Context) error {
	return api.Execute(e, c.productService.UpdateProduct)
}

func (c *productController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productService.DeleteProduct(ctx, id)
	})
}
