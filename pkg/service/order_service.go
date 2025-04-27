package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/enum"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"time"
)

type OrderService interface {
	ListOrders(ctx context.Context, filter dto.OrderListFilter) (api.Response, error)
	GetOrderById(ctx context.Context, id string) (api.Response, error)
	CreateOrder(ctx context.Context, req dto.OrderCreate) (api.Response, error)
	UpdateOrder(ctx context.Context, req dto.OrderUpdate) (api.Response, error)
	DeleteOrder(ctx context.Context, id string) (api.Response, error)
	UpdateOrderStatus(ctx context.Context, id string, status enum.OrderStatus) (api.Response, error)
}

type orderService struct {
	db *supabase.Client
}

func NewOrderService(di do.Injector) (OrderService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to init OrderService: %w", err)
	}
	return &orderService{db: db}, nil
}

/* ------------------------------------------------------------------
   Helpers
   ------------------------------------------------------------------*/

func (s *orderService) countOrders(ctx context.Context, filter dto.OrderListFilter) (int, error) {
	q := s.db.DB.From("orders").Select("id").IsNull("deleted_at")
	if filter.Status != "" {
		q = q.Eq("status", string(filter.Status))
	}

	var tmp []struct{}
	if err := q.Execute(&tmp); err != nil {
		return 0, fmt.Errorf("failed to count orders: %w", err)
	}
	return len(tmp), nil
}

// make sure every option referenced in an order still exists
func (s *orderService) validateOptionIds(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("order must contain at least one product option")
	}
	var found []struct {
		Id string `json:"id"`
	}
	if err := s.db.DB.
		From("product_options").
		Select("id").
		In("id", ids).
		IsNull("deleted_at").
		Execute(&found); err != nil {
		return fmt.Errorf("failed to validate product options: %w", err)
	}
	if len(found) != len(ids) {
		return fmt.Errorf("one or more product options not found or deleted")
	}
	return nil
}

/* ------------------------------------------------------------------
   List
   ------------------------------------------------------------------*/

func (s *orderService) ListOrders(ctx context.Context, filter dto.OrderListFilter) (api.Response, error) {
	total, err := s.countOrders(ctx, filter)
	if err != nil {
		return nil, err
	}

	var orders []dto.Order
	q := s.db.DB.
		From("orders").
		Select("id,created_at,updated_at,user_id,address,status").
		LimitWithOffset(int(filter.Limit), int((filter.Page-1)*filter.Limit)).
		IsNull("deleted_at")

	if filter.Status != "" {
		q = q.Eq("status", string(filter.Status))
	}

	if err := q.Execute(&orders); err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	filter.Total = int64(total)
	return api.SuccessPagination(orders, &filter.Pagination), nil
}

/* ------------------------------------------------------------------
   Get by id (base order + order_details)
   ------------------------------------------------------------------*/

func (s *orderService) GetOrderById(ctx context.Context, id string) (api.Response, error) {
	// 1) base order ---------------------------------------------------------
	var orders []dto.Order
	if err := s.db.DB.
		From("orders").
		Select("id,created_at,updated_at,user_id,address,status,metadata").
		Eq("id", id).
		IsNull("deleted_at").
		Execute(&orders); err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	// 2) details ------------------------------------------------------------
	var details []dto.OrderDetail
	if err := s.db.DB.
		From("order_details").
		Select("id,order_id,product_option_id,quantity,discount,discount_type,metadata").
		Eq("order_id", id).
		IsNull("deleted_at").
		Execute(&details); err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	resp := dto.OrderDetailResponse{
		Order:       orders[0],
		OrderDetail: details,
	}
	return api.Success(resp), nil
}

/* ------------------------------------------------------------------
   Create
   ------------------------------------------------------------------*/

func (s *orderService) CreateOrder(ctx context.Context, req dto.OrderCreate) (api.Response, error) {
	// 1) validate options exist
	var optionIds []string
	for _, od := range req.OrderDetails {
		optionIds = append(optionIds, od.ProductOptionId)
	}
	if err := s.validateOptionIds(optionIds); err != nil {
		return nil, err
	}

	// 2) insert into orders --------------------------------------------------
	orderBody := map[string]interface{}{
		"user_id":  req.UserId,
		"address":  req.Address,
		"status":   req.Status,
		"metadata": req.Metadata,
	}
	var createdOrders []dto.Order
	if err := s.db.DB.From("orders").Insert(orderBody).Execute(&createdOrders); err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}
	orderId := createdOrders[0].Id

	// 3) insert order_details -----------------------------------------------
	var odRows []map[string]interface{}
	for _, od := range req.OrderDetails {
		odRows = append(odRows, map[string]interface{}{
			"order_id":          orderId,
			"product_option_id": od.ProductOptionId,
			"quantity":          od.Quantity,
			"discount":          od.Discount,
			"discount_type":     od.DiscountType,
			"metadata":          od.Metadata,
		})
	}

	if err := s.db.DB.From("order_details").Insert(odRows).Execute(nil); err != nil {
		// best-effort rollback
		_ = s.db.DB.From("orders").Delete().Eq("id", orderId).Execute(nil)
		return nil, fmt.Errorf("failed to insert order details: %v", err)
	}

	return s.GetOrderById(ctx, orderId)
}

/* ------------------------------------------------------------------
   Update
   ------------------------------------------------------------------*/

func (s *orderService) UpdateOrder(ctx context.Context, req dto.OrderUpdate) (api.Response, error) {
	// 1) validate status
	if req.Status != enum.Pending &&
		req.Status != enum.Complete &&
		req.Status != enum.Cancelled {
		return nil, fmt.Errorf("invalid status")
	}

	// 2) validate options
	var optionIds []string
	for _, od := range req.OrderDetails {
		optionIds = append(optionIds, od.ProductOptionId)
	}
	if err := s.validateOptionIds(optionIds); err != nil {
		return nil, err
	}

	// 3) update order
	updateBody := map[string]interface{}{
		"user_id":    req.UserId,
		"address":    req.Address,
		"status":     req.Status,
		"metadata":   req.Metadata,
		"updated_at": time.Now(),
	}
	if err := s.db.DB.
		From("orders").
		Update(updateBody).
		Eq("id", req.Id).
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to update order: %v", err)
	}

	// 4) replace order_details – simplest approach: soft-delete old rows & re-insert
	if err := s.db.DB.
		From("order_details").
		Update(map[string]interface{}{"deleted_at": time.Now()}).
		Eq("order_id", req.Id).
		IsNull("deleted_at").
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to clear old order details: %v", err)
	}
	var newRows []map[string]interface{}
	for _, od := range req.OrderDetails {
		newRows = append(newRows, map[string]interface{}{
			"order_id":          req.Id,
			"product_option_id": od.ProductOptionId,
			"quantity":          od.Quantity,
			"discount":          od.Discount,
			"discount_type":     od.DiscountType,
			"metadata":          od.Metadata,
		})
	}
	if err := s.db.DB.From("order_details").Insert(newRows).Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to add new order details: %v", err)
	}

	return s.GetOrderById(ctx, req.Id)
}

/* ------------------------------------------------------------------
   Soft-delete
   ------------------------------------------------------------------*/

func (s *orderService) DeleteOrder(ctx context.Context, id string) (api.Response, error) {
	now := time.Now()
	if err := s.db.DB.
		From("orders").
		Update(map[string]interface{}{"deleted_at": now}).
		Eq("id", id).
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to delete order: %v", err)
	}
	if err := s.db.DB.
		From("order_details").
		Update(map[string]interface{}{"deleted_at": now}).
		Eq("order_id", id).
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to delete order details: %v", err)
	}
	return api.Success("Order deleted successfully"), nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status enum.OrderStatus) (api.Response, error) {
	// Validate the incoming status
	if status != enum.Pending &&
		status != enum.Complete &&
		status != enum.Cancelled {
		return nil, fmt.Errorf("invalid status")
	}

	updateBody := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if err := s.db.DB.
		From("orders").
		Update(updateBody).
		Eq("id", id).
		IsNull("deleted_at").
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to update order status: %v", err)
	}

	return api.Success("Order status updated successfully"), nil
}
