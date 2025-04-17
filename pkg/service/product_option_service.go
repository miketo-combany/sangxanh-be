package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

type ProductOptionService interface {
	ListProductOptions(ctx context.Context, productId string) (api.Response, error)
	CreateProductOption(ctx context.Context, req dto.ProductOptionCreate) (api.Response, error)
	UpdateProductOption(ctx context.Context, req dto.ProductOptionUpdate) (api.Response, error)
	DeleteProductOption(ctx context.Context, id string) (api.Response, error)
}

type productOptionService struct {
	db *supabase.Client
}

func NewProductOptionService(di do.Injector) (ProductOptionService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ProductOptionService: %w", err)
	}
	return &productOptionService{db: db}, nil
}

func (s *productOptionService) ListProductOptions(ctx context.Context, productId string) (api.Response, error) {
	var options []dto.ProductOption
	err := s.db.DB.
		From("product_options").
		Select("id,name,product_id,price,detail,metadata,created_at,updated_at").
		Eq("product_id", productId).
		IsNull("deleted_at").
		Execute(&options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product options: %v", err)
	}
	return api.Success(options), nil
}

func (s *productOptionService) CreateProductOption(ctx context.Context, req dto.ProductOptionCreate) (api.Response, error) {
	// Validate product existence
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	var created []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Insert(req).
		Execute(&created); err != nil {
		return nil, fmt.Errorf("failed to create product option: %v", err)
	}
	return api.Success(created), nil
}

func (s *productOptionService) UpdateProductOption(ctx context.Context, req dto.ProductOptionUpdate) (api.Response, error) {
	updateData := map[string]interface{}{
		"name":       req.Name,
		"price":      req.Price,
		"detail":     req.Detail,
		"metadata":   req.Metadata,
		"updated_at": time.Now(),
	}

	var updated []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Update(updateData).
		Eq("id", req.Id).
		Execute(&updated); err != nil {
		return nil, fmt.Errorf("failed to update product option: %v", err)
	}

	return api.Success(updated), nil
}

func (s *productOptionService) DeleteProductOption(ctx context.Context, id string) (api.Response, error) {
	updateData := map[string]interface{}{
		"deleted_at": time.Now(),
	}

	var updated []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Update(updateData).
		Eq("id", id).
		Execute(&updated); err != nil {
		return nil, fmt.Errorf("failed to delete product option: %v", err)
	}

	return api.Success("Product option deleted successfully"), nil
}

func (s *productOptionService) validProduct(id string) error {
	var product []dto.Product
	err := s.db.DB.
		From("products").
		Select("id").
		Eq("id", id).
		IsNull("deleted_at").
		Execute(&product)
	if err != nil {
		return fmt.Errorf("failed to find product: %v", err)
	}
	if len(product) == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
