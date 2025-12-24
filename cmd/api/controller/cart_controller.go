package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/service"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type cartController struct {
	cartService    service.CartService
	authMiddleware echo.MiddlewareFunc
}

func NewCartController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &cartController{
		cartService:    do.MustInvoke[service.CartService](di),
		authMiddleware: auth,
	}, nil
}

func (c *cartController) Register(g *echo.Group) {
	g = g.Group("/cart")
	g.GET("", c.List, c.authMiddleware)             // List all carts for the current user
	g.POST("/create", c.Create, c.authMiddleware)   // Create a new cart
	g.PUT("/update", c.Update, c.authMiddleware)    // Update cart quantity
	g.DELETE("/delete", c.Delete, c.authMiddleware) // Delete a cart
}

// List godoc
// @Summary List user's cart items
// @Description Get all cart items for the authenticated user
// @Tags Cart
// @Accept json
// @Produce json
// @Success 200 {array} dto.CartResponse "List of cart items"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /cart [get]
func (c *cartController) List(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.cartService.GetCartsByUserID(ctx)
	})
}

// Create godoc
// @Summary Add item to cart
// @Description Add a product option to the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body dto.CartCreateRequest true "Cart item data"
// @Success 201 {object} dto.CartResponse "Created cart item"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /cart/create [post]
func (c *cartController) Create(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, req dto.CartCreateRequest) (api.Response, error) {
		return c.cartService.CreateCart(ctx, req)
	})
}

// Update godoc
// @Summary Update cart item quantity
// @Description Update the quantity of a cart item
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body dto.CartUpdate true "Cart update data"
// @Success 200 {object} dto.CartResponse "Updated cart item"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Cart item not found"
// @Security BearerAuth
// @Router /cart/update [put]
func (c *cartController) Update(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, req dto.CartUpdate) (api.Response, error) {
		return c.cartService.UpdateCart(ctx, req)
	})
}

// Delete godoc
// @Summary Remove item from cart
// @Description Delete a cart item by ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param id query string true "Cart item ID"
// @Success 200 {object} map[string]interface{} "Cart item deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Cart item not found"
// @Security BearerAuth
// @Router /cart/delete [delete]
func (c *cartController) Delete(e echo.Context) error {
	id := e.QueryParam("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.cartService.DeleteCart(ctx, id)
	})
}
