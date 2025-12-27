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
	authMiddleware        echo.MiddlewareFunc
}

func NewProductVariantController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &productVariantController{
		productVariantService: do.MustInvoke[service.ProductVariantService](di),
		authMiddleware:        auth,
	}, nil
}

func (c *productVariantController) Register(g *echo.Group) {
	g = g.Group("/product-variant")
	g.GET("", c.List)
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update)
	g.DELETE("/delete", c.Delete)
	g.POST("/create-bulk", c.CreateBulk)
	g.PUT("/update-bulk", c.UpdateBulk)
}

// List godoc
// @Summary List product variants
// @Description Get all variants for a specific product
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param productId query string true "Product ID"
// @Success 200 {array} dto.ProductVariantResponse "List of product variants"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-variant [get]
func (c *productVariantController) List(e echo.Context) error {
	id := e.QueryParam("productId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productVariantService.ListProductVariants(ctx, id)
	})
}

// Create godoc
// @Summary Create a product variant
// @Description Create a new variant for a product
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param request body dto.ProductVariantCreate true "Product variant data"
// @Success 201 {object} dto.ProductVariantResponse "Created product variant"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-variant/create [post]
func (c *productVariantController) Create(e echo.Context) error {
	return api.Execute(e, c.productVariantService.CreateProductVariant)
}

// Update godoc
// @Summary Update a product variant
// @Description Update an existing product variant
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param request body dto.ProductVariantUpdate true "Product variant update data"
// @Success 200 {object} dto.ProductVariantResponse "Updated product variant"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Product variant not found"
// @Router /product-variant/update [put]
func (c *productVariantController) Update(e echo.Context) error {
	return api.Execute(e, c.productVariantService.UpdateProductVariant)
}

// Delete godoc
// @Summary Delete a product variant
// @Description Delete a product variant by ID
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param id query string true "Product variant ID"
// @Success 200 {object} map[string]interface{} "Product variant deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Product variant not found"
// @Router /product-variant/delete [delete]
func (c *productVariantController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productVariantService.DeleteProductVariant(ctx, id)
	})
}

// CreateBulk godoc
// @Summary Create multiple product variants
// @Description Create multiple variants for a product in one request
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param request body dto.ProductVariantCreateBulk true "Bulk product variants data"
// @Success 201 {array} dto.ProductVariantResponse "Created product variants"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-variant/create-bulk [post]
func (c *productVariantController) CreateBulk(e echo.Context) error {
	return api.Execute(e, c.productVariantService.CreateBulkProductVariant)
}

// UpdateBulk godoc
// @Summary Update multiple product variants
// @Description Update multiple variants for a product in one request
// @Tags Product Variants
// @Accept json
// @Produce json
// @Param request body dto.ProductVariantUpdateBulk true "Bulk product variants update data"
// @Success 200 {array} dto.ProductVariantResponse "Updated product variants"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-variant/update-bulk [put]
func (c *productVariantController) UpdateBulk(e echo.Context) error {
	return api.Execute(e, c.productVariantService.UpdateBulkProductVariant)
}
