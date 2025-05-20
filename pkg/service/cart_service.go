package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"time"
)

type CartService interface {
	CreateCart(ctx context.Context, req dto.CartCreateRequest, userID string) (api.Response, error)
	GetCartsByUserID(ctx context.Context, userID string) (api.Response, error)
	UpdateCart(ctx context.Context, req dto.CartUpdate) (api.Response, error)
	DeleteCart(ctx context.Context, id string) (api.Response, error)
}

type cartService struct {
	db *supabase.Client
}

func NewCartService(di do.Injector) (CartService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CartService: %w", err)
	}
	return &cartService{db: db}, nil
}

func (s *cartService) CreateCart(ctx context.Context, req dto.CartCreateRequest, userID string) (api.Response, error) {
	var created []dto.Cart

	if err := s.db.DB.
		From("carts").
		Insert(req).
		Execute(&created); err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}
	// Return the created cart
	return api.Success(created[0]), nil
}

func (s *cartService) GetCartsByUserID(ctx context.Context, userID string) (api.Response, error) {
	// Query to get carts for a specific user
	var carts []dto.Cart
	if err := s.db.DB.
		From("carts").
		Select("id, user_id, product_option_id, quantity, created_at, updated_at").
		Eq("user_id", userID).
		IsNull("deleted_at").
		Execute(&carts); err != nil {
		return nil, fmt.Errorf("failed to fetch carts for user %s: %w", userID, err)
	}

	// Now query the product_options table to get ProductOption data for each cart
	var cartResponses []dto.CartResponse
	for _, cart := range carts {
		var productOption dto.ProductOption
		if err := s.db.DB.
			From("product_options").
			Select("id, name, price, detail, created_at, updated_at").
			Eq("id", cart.ProductOptionID).
			IsNull("deleted_at").
			Execute(&productOption); err != nil {
			return nil, fmt.Errorf("failed to fetch product option for cart %s: %w", cart.ID, err)
		}

		// Combine cart and product option in CartResponse
		cartResponse := dto.CartResponse{
			Cart:          cart,
			ProductOption: productOption,
		}
		cartResponses = append(cartResponses, cartResponse)
	}

	// Return the combined response with cart and product option
	return api.Success(cartResponses), nil
}

func (s *cartService) UpdateCart(ctx context.Context, req dto.CartUpdate) (api.Response, error) {
	updateData := map[string]interface{}{
		"quantity":   req.Quantity,
		"updated_at": time.Now(),
	}

	var updated []dto.Cart
	if err := s.db.DB.
		From("carts").
		Update(updateData).
		Eq("id", req.ID).
		Execute(&updated); err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	// Fetch the updated cart and return with ProductOption
	var productOption dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Select("id, name, price, detail, created_at, updated_at").
		Eq("id", updated[0].ProductOptionID).
		IsNull("deleted_at").
		Execute(&productOption); err != nil {
		return nil, fmt.Errorf("failed to fetch product option for cart %s: %w", updated[0].ID, err)
	}

	cartResponse := dto.CartResponse{
		Cart:          updated[0],
		ProductOption: productOption,
	}

	return api.Success(cartResponse), nil
}

func (s *cartService) DeleteCart(ctx context.Context, id string) (api.Response, error) {
	updateData := map[string]interface{}{
		"deleted_at": time.Now(),
	}

	var updated []dto.Cart
	if err := s.db.DB.
		From("carts").
		Update(updateData).
		Eq("id", id).
		Execute(&updated); err != nil {
		return nil, fmt.Errorf("failed to delete cart: %w", err)
	}

	// Return a success message
	return api.Success("Cart deleted successfully"), nil
}
