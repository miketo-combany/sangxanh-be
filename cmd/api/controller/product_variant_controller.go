package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type productVariantController struct {
	productVariantService service.ProductVariantService
}

func NewProductVariantController(di do.Injector) (api.Controller, error) {
	return &productVariantController{
		productVariantService: do.MustInvoke[service.ProductVariantService](di),
	}, nil
}

func (c *productVariantController) Register(g *echo.Group) {
	g = g.Group("/product-variant")
	g.GET("", c.List)
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update)
	g.DELETE("/delete", c.Delete)
}

func (c *productVariantController) List(e echo.Context) error {
	id := e.QueryParam("productId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productVariantService.ListProductVariants(ctx, id)
	})
}

func (c *productVariantController) Create(e echo.Context) error {
	return api.Execute(e, c.productVariantService.CreateProductVariant)
}

func (c *productVariantController) Update(e echo.Context) error {
	return api.Execute(e, c.productVariantService.UpdateProductVariant)
}

func (c *productVariantController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productVariantService.DeleteProductVariant(ctx, id)
	})
}
