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

type productController struct {
	productService service.ProductService
	authMiddleware echo.MiddlewareFunc
}

func NewProductController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &productController{
		productService: do.MustInvoke[service.ProductService](di),
		authMiddleware: auth,
	}, nil
}

func (c *productController) Register(g *echo.Group) {
	g = g.Group("/product")
	g.GET("", c.List)
	g.POST("/create", c.Create, c.authMiddleware, middleware.RequireRoles("admin"))
	g.PUT("/update", c.Update, c.authMiddleware, middleware.RequireRoles("admin"))
	g.DELETE("/delete", c.Delete, c.authMiddleware, middleware.RequireRoles("admin"))
	g.GET("/:id", c.GetById)
}

func (c *productController) List(e echo.Context) error {
	return api.Execute[dto.ProductFilter](e, func(
		ctx context.Context,
		req dto.ProductFilter, // ← everything is inside req now
	) (api.Response, error) {
		return c.productService.ListProducts(ctx, req)
	})
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

func (c *productController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productService.GetProductById(ctx, id)
	})
}
