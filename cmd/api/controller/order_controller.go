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
	g.POST("/create", c.Create)
	g.PUT("/update", c.Update, c.authMiddleware)
	g.DELETE("/delete", c.Delete, c.authMiddleware)
	g.PUT("/update-status", c.UpdateStatus, c.authMiddleware)
}

// List godoc
// @Summary List orders
// @Description Get a paginated list of orders with optional filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Filter by order status"
// @Param user_id query string false "Filter by user ID"
// @Success 200 {object} map[string]interface{} "Paginated list of orders"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /order [get]
func (c *orderController) List(e echo.Context) error {
	return api.Execute[dto.OrderListFilter](e, func(ctx context.Context, req dto.OrderListFilter) (api.Response, error) {
		return c.orderService.ListOrders(ctx, req)
	})
}

// GetById godoc
// @Summary Get order by ID
// @Description Retrieve a specific order by its ID with order details
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderDetailResponse "Order details"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Router /order/{id} [get]
func (c *orderController) GetById(e echo.Context) error {
	id := e.Param("id")
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.orderService.GetOrderById(ctx, id)
	})
}

// Create godoc
// @Summary Create a new order
// @Description Create a new order with order details
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.OrderCreate true "Order data"
// @Success 201 {object} dto.OrderDetailResponse "Created order"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /order/create [post]
func (c *orderController) Create(e echo.Context) error {
	return api.Execute(e, c.orderService.CreateOrder)
}

// Update godoc
// @Summary Update an order
// @Description Update an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.OrderUpdate true "Order update data"
// @Success 200 {object} dto.OrderDetailResponse "Updated order"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security BearerAuth
// @Router /order/update [put]
func (c *orderController) Update(e echo.Context) error {
	return api.Execute(e, c.orderService.UpdateOrder)
}

// Delete godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param orderId query string true "Order ID"
// @Success 200 {object} map[string]interface{} "Order deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security BearerAuth
// @Router /order/delete [delete]
func (c *orderController) Delete(e echo.Context) error {
	id := e.QueryParam("orderId") // Consistent with your DeleteCategory
	return api.Execute(e, func(ctx context.Context, _ struct{}) (api.Response, error) {
		return c.orderService.DeleteOrder(ctx, id)
	})
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body object{order_id=string,status=string} true "Order status update"
// @Success 200 {object} dto.OrderDetailResponse "Updated order"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security BearerAuth
// @Router /order/update-status [put]
func (c *orderController) UpdateStatus(e echo.Context) error {
	type Req struct {
		OrderId string           `json:"order_id" validate:"required"`
		Status  enum.OrderStatus `json:"status" validate:"required"`
	}

	return api.Execute[Req](e, func(ctx context.Context, req Req) (api.Response, error) {
		return c.orderService.UpdateOrderStatus(ctx, req.OrderId, req.Status)
	})
}
