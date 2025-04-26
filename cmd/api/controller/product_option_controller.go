package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type productOptionController struct {
	productOption service.ProductOptionService
}

func NewProductOptionController(di do.Injector) (api.Controller, error) {
	return &productOptionController{
		productOption: do.MustInvoke[service.ProductOptionService](di),
	}, nil
}

func (c *productOptionController) Register(g *echo.Group) {
	g = g.Group("/product-option")
	g.GET("", c.List)
	g.POST("/create", c.Create)
	g.POST("/create-bulk", c.CreateBulk)
	g.PUT("/update", c.Update)
	g.PUT("/update-bulk", c.UpdateBulk)
	g.DELETE("/delete", c.Delete)
}

func (c *productOptionController) List(e echo.Context) error {
	id := e.QueryParam("productId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productOption.ListProductOptions(ctx, id)
	})
}

func (c *productOptionController) Create(e echo.Context) error {
	return api.Execute(e, c.productOption.CreateProductOption)
}

func (c *productOptionController) Update(e echo.Context) error {
	return api.Execute(e, c.productOption.UpdateProductOption)
}

func (c *productOptionController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productOption.DeleteProductOption(ctx, id)
	})
}

func (c *productOptionController) CreateBulk(e echo.Context) error {
	return api.Execute(e, c.productOption.CreateBulkProductOption)
}

func (c *productOptionController) UpdateBulk(e echo.Context) error {
	return api.Execute(e, c.productOption.UpdateBulkProductOption)
}
