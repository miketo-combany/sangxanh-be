package controller

import (
	"SangXanh/cmd/api/middleware"
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/service"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type shopController struct {
	shopService service.ShopService
	authMiddleware  echo.MiddlewareFunc
}

func NewShopController(di do.Injector, authMiddleware echo.MiddlewareFunc) (api.Controller, error) {
	return &shopController{
		shopService: do.MustInvoke[service.ShopService](di),
		authMiddleware:  authMiddleware,
	}, nil
}

func (c *shopController) Register(g *echo.Group) {
	g = g.Group("/shop-a")
	g.GET("", c.GetShop)
	g.PUT("", c.UpdateShop, c.authMiddleware, middleware.RequireRoles("admin"))
}

// GetShop godoc
// @Summary Get shop information
// @Description Retrieve the shop information (only one shop record exists)
// @Tags Shop
// @Accept json
// @Produce json
// @Success 200 {object} dto.ShopResponse "Shop information"
// @Failure 404 {object} map[string]interface{} "Shop not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /shop [get]
func (c *shopController) GetShop(e echo.Context) error {
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.shopService.GetShop(ctx)
	})
}

// UpdateShop godoc
// @Summary Update shop information
// @Description Update the shop information (admin only, updates the first record)
// @Tags Shop
// @Accept json
// @Produce json
// @Param request body dto.ShopUpdate true "Shop update data"
// @Success 200 {object} dto.ShopResponse "Updated shop information"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin only"
// @Failure 404 {object} map[string]interface{} "Shop not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /shop [put]
func (c *shopController) UpdateShop(e echo.Context) error {
	return api.Execute(e, c.shopService.UpdateShop)
}
