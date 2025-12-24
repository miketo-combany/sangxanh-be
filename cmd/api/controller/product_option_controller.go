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
	auth          echo.MiddlewareFunc
}

func NewProductOptionController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &productOptionController{
		productOption: do.MustInvoke[service.ProductOptionService](di),
		auth:          auth,
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

// List godoc
// @Summary List product options
// @Description Get all options for a specific product
// @Tags Product Options
// @Accept json
// @Produce json
// @Param productId query string true "Product ID"
// @Success 200 {array} dto.ProductOptionResponse "List of product options"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-option [get]
func (c *productOptionController) List(e echo.Context) error {
	id := e.QueryParam("productId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productOption.ListProductOptions(ctx, id)
	})
}

// Create godoc
// @Summary Create a product option
// @Description Create a new option for a product
// @Tags Product Options
// @Accept json
// @Produce json
// @Param request body dto.ProductOptionCreate true "Product option data"
// @Success 201 {object} dto.ProductOptionResponse "Created product option"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-option/create [post]
func (c *productOptionController) Create(e echo.Context) error {
	return api.Execute(e, c.productOption.CreateProductOption)
}

// Update godoc
// @Summary Update a product option
// @Description Update an existing product option
// @Tags Product Options
// @Accept json
// @Produce json
// @Param request body dto.ProductOptionUpdate true "Product option update data"
// @Success 200 {object} dto.ProductOptionResponse "Updated product option"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Product option not found"
// @Router /product-option/update [put]
func (c *productOptionController) Update(e echo.Context) error {
	return api.Execute(e, c.productOption.UpdateProductOption)
}

// Delete godoc
// @Summary Delete a product option
// @Description Delete a product option by ID
// @Tags Product Options
// @Accept json
// @Produce json
// @Param id query string true "Product option ID"
// @Success 200 {object} map[string]interface{} "Product option deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Product option not found"
// @Router /product-option/delete [delete]
func (c *productOptionController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productOption.DeleteProductOption(ctx, id)
	})
}

// CreateBulk godoc
// @Summary Create multiple product options
// @Description Create multiple options for a product in one request
// @Tags Product Options
// @Accept json
// @Produce json
// @Param request body dto.ProductOptionCreateBulk true "Bulk product options data"
// @Success 201 {array} dto.ProductOptionResponse "Created product options"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-option/create-bulk [post]
func (c *productOptionController) CreateBulk(e echo.Context) error {
	return api.Execute(e, c.productOption.CreateBulkProductOption)
}

// UpdateBulk godoc
// @Summary Update multiple product options
// @Description Update multiple options for a product in one request
// @Tags Product Options
// @Accept json
// @Produce json
// @Param request body dto.ProductOptionBulkUpdate true "Bulk product options update data"
// @Success 200 {array} dto.ProductOptionResponse "Updated product options"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product-option/update-bulk [put]
func (c *productOptionController) UpdateBulk(e echo.Context) error {
	return api.Execute(e, c.productOption.UpdateBulkProductOption)
}
