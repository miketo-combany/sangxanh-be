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

// List godoc
// @Summary List all products
// @Description Get a paginated list of products with optional filters
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param name query string false "Filter by product name"
// @Param category_id query string false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {object} map[string]interface{} "Paginated list of products"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /product [get]
func (c *productController) List(e echo.Context) error {
	return api.Execute[dto.ProductFilter](e, func(
		ctx context.Context,
		req dto.ProductFilter, // ← everything is inside req now
	) (api.Response, error) {
		return c.productService.ListProducts(ctx, req)
	})
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.ProductCreated true "Product data"
// @Success 201 {object} dto.ProductResponse "Created product"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Security BearerAuth
// @Router /product/create [post]
func (c *productController) Create(e echo.Context) error {
	return api.Execute(e, c.productService.CreateProduct)
}

// Update godoc
// @Summary Update a product
// @Description Update an existing product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.ProductUpdated true "Product update data"
// @Success 200 {object} dto.ProductResponse "Updated product"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Security BearerAuth
// @Router /product/update [put]
func (c *productController) Update(e echo.Context) error {
	return api.Execute(e, c.productService.UpdateProduct)
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param id query string true "Product ID"
// @Success 200 {object} map[string]interface{} "Product deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Security BearerAuth
// @Router /product/delete [delete]
func (c *productController) Delete(e echo.Context) error {
	id := e.QueryParam("id") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productService.DeleteProduct(ctx, id)
	})
}

// GetById godoc
// @Summary Get product by ID
// @Description Retrieve a specific product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse "Product details"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /product/{id} [get]
func (c *productController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.productService.GetProductById(ctx, id)
	})
}
