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
	g.GET("/list-header", c.ListHeader)
}

// List godoc
// @Summary List all categories
// @Description Get a paginated list of categories with optional filters
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param name query string false "Filter by category name"
// @Success 200 {object} map[string]interface{} "Paginated list of categories"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /category [get]
func (c *categoryController) List(e echo.Context) error {
	return api.Execute[dto.ListCategory](e, func(ctx context.Context, req dto.ListCategory) (api.Response, error) {
		return c.categoryService.ListCategories(ctx, req)
	})
}

// GetById godoc
// @Summary Get category by ID
// @Description Retrieve a specific category by its ID with subcategories
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse "Category details"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Router /category/{id} [get]
func (c *categoryController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.ListCategoryById(ctx, id)
	})
}

// Create godoc
// @Summary Create a new category
// @Description Create a new category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CategoryCreate true "Category data"
// @Success 201 {object} dto.CategoryResponse "Created category"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Security BearerAuth
// @Router /category/create [post]
func (c *categoryController) Create(e echo.Context) error {
	return api.Execute(e, c.categoryService.CreateCategory)
}

// Update godoc
// @Summary Update a category
// @Description Update an existing category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CategoryUpdate true "Category update data"
// @Success 200 {object} dto.CategoryResponse "Updated category"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /category/update [put]
func (c *categoryController) Update(e echo.Context) error {
	return api.Execute(e, c.categoryService.UpdateCategory)
}

// Delete godoc
// @Summary Delete a category
// @Description Delete a category by ID (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Param categoryId query string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin only"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Security BearerAuth
// @Router /category/delete [delete]
func (c *categoryController) Delete(e echo.Context) error {
	id := e.QueryParam("categoryId") // Get "id" parameter from URL
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.DeleteCategory(ctx, id)
	})
}

// ListHeader godoc
// @Summary List header categories
// @Description Get categories marked for display in header
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoryListResponse "List of header categories"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /category/list-header [get]
func (c *categoryController) ListHeader(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.categoryService.ListHeaderCategories(ctx)
	})
}
