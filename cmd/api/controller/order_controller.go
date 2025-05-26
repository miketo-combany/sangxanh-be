package controller

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/enum"
	"SangXanh/pkg/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type orderController struct {
	orderService   service.OrderService
	authMiddleware echo.MiddlewareFunc
}

func NewOrderController(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error) {
	return &orderController{
		orderService:   do.MustInvoke[service.OrderService](di),
		authMiddleware: auth,
	}, nil
}

func (c *orderController) Register(g *echo.Group) {
	g = g.Group("/order")
	g.GET("", c.List, c.authMiddleware)
	g.GET("/:id", c.GetById)
	g.POST("/create", c.Create, c.authMiddleware)
	g.PUT("/update", c.Update, c.authMiddleware)
	g.DELETE("/delete", c.Delete, c.authMiddleware)
	g.PUT("/update-status", c.UpdateStatus, c.authMiddleware)
}

func (c *orderController) List(e echo.Context) error {
	return api.Execute[dto.OrderListFilter](e, func(ctx context.Context, req dto.OrderListFilter) (api.Response, error) {
		return c.orderService.ListOrders(ctx, req)
	})
}

func (c *orderController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.orderService.GetOrderById(ctx, id)
	})
}

func (c *orderController) Create(e echo.Context) error {
	return api.Execute(e, c.orderService.CreateOrder)
}

func (c *orderController) Update(e echo.Context) error {
	return api.Execute(e, c.orderService.UpdateOrder)
}

func (c *orderController) Delete(e echo.Context) error {
	id := e.QueryParam("orderId") // Consistent with your DeleteCategory
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.orderService.DeleteOrder(ctx, id)
	})
}

func (c *orderController) UpdateStatus(e echo.Context) error {
	type Req struct {
		OrderId string           `json:"order_id" validate:"required"`
		Status  enum.OrderStatus `json:"status" validate:"required"`
	}

	return api.Execute[Req](e, func(ctx context.Context, req Req) (api.Response, error) {
		return c.orderService.UpdateOrderStatus(ctx, req.OrderId, req.Status)
	})
}
