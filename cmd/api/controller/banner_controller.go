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

type bannerController struct {
	bannerService  service.BannerService
	authMiddleware echo.MiddlewareFunc
}

func NewBannerController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &bannerController{
		bannerService:  do.MustInvoke[service.BannerService](di),
		authMiddleware: auth,
	}, nil
}

func (c *bannerController) Register(g *echo.Group) {
	g = g.Group("/banner")
	g.GET("", c.List)
	g.GET("/:id", c.GetById)
	g.POST("/create", c.Create, c.authMiddleware, middleware.RequireRoles("admin"))
	g.PUT("/update", c.Update, c.authMiddleware, middleware.RequireRoles("admin"))
	g.DELETE("/delete/:id", c.Delete, c.authMiddleware, middleware.RequireRoles("admin"))
}

// List godoc
// @Summary List all banners
// @Description Get a list of banners with optional filters
// @Tags Banners
// @Accept json
// @Produce json
// @Param slot query string false "Filter by slot (first, second, third)"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} api.Response{data=[]dto.Banner}
// @Failure 500 {object} api.Response
// @Router /banner [get]
func (c *bannerController) List(e echo.Context) error {
	return api.Execute[dto.BannerFilter](e, func(ctx context.Context, req dto.BannerFilter) (api.Response, error) {
		return c.bannerService.ListBanners(ctx, req)
	})
}

// GetById godoc
// @Summary Get a banner by ID
// @Description Get detailed information about a specific banner
// @Tags Banners
// @Accept json
// @Produce json
// @Param id path string true "Banner ID"
// @Success 200 {object} api.Response{data=dto.Banner}
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /banner/{id} [get]
func (c *bannerController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.bannerService.GetBannerById(ctx, id)
	})
}

// Create godoc
// @Summary Create a new banner
// @Description Create a new banner (Admin only)
// @Tags Banners
// @Accept json
// @Produce json
// @Param banner body dto.BannerCreate true "Banner creation data"
// @Success 201 {object} api.Response{data=dto.Banner}
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 403 {object} api.Response
// @Failure 500 {object} api.Response
// @Security Bearer
// @Router /banner/create [post]
func (c *bannerController) Create(e echo.Context) error {
	return api.Execute[dto.BannerCreate](e, func(ctx context.Context, req dto.BannerCreate) (api.Response, error) {
		return c.bannerService.CreateBanner(ctx, req)
	})
}

// Update godoc
// @Summary Update a banner
// @Description Update an existing banner (Admin only)
// @Tags Banners
// @Accept json
// @Produce json
// @Param banner body dto.BannerUpdate true "Banner update data"
// @Success 200 {object} api.Response{data=dto.Banner}
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 403 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Security Bearer
// @Router /banner/update [put]
func (c *bannerController) Update(e echo.Context) error {
	return api.Execute[dto.BannerUpdate](e, func(ctx context.Context, req dto.BannerUpdate) (api.Response, error) {
		return c.bannerService.UpdateBanner(ctx, req)
	})
}

// Delete godoc
// @Summary Delete a banner
// @Description Delete a banner by ID (Admin only)
// @Tags Banners
// @Accept json
// @Produce json
// @Param id path string true "Banner ID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 401 {object} api.Response
// @Failure 403 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Security Bearer
// @Router /banner/delete/{id} [delete]
func (c *bannerController) Delete(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.bannerService.DeleteBanner(ctx, id)
	})
}
