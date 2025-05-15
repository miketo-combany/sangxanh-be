package controller

import (
	"SangXanh/pkg/common/api"
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
	g.GET("", c.List)             // List all carts for the current user
	g.POST("/create", c.Create)   // Create a new cart
	g.PUT("/update", c.Update)    // Update cart quantity
	g.DELETE("/delete", c.Delete) // Delete a cart
}

func (c *cartController) List(e echo.Context) error {
	userID := e.Get("userID").(string) // Get user ID from the access token
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.cartService.GetCartsByUserID(ctx, userID)
	})
}

func (c *cartController) Create(e echo.Context) error {
	return api.Execute(e, c.cartService.CreateCart)
}

func (c *cartController) Update(e echo.Context) error {
	return api.Execute(e, c.cartService.UpdateCart)
}

func (c *cartController) Delete(e echo.Context) error {
	id := e.QueryParam("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.cartService.DeleteCart(ctx, id)
	})
}
