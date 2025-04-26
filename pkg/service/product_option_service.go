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
	CreateBulkProductOption(ctx context.Context, req dto.ProductOptionCreateBulk) (api.Response, error)
	UpdateBulkProductOption(ctx context.Context, req dto.ProductOptionBulkUpdate) (api.Response, error)
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
	return api.Success(options[0]), nil
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
	return api.Success(created[0]), nil
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

	return api.Success(updated[0]), nil
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

func (s *productOptionService) CreateBulkProductOption(ctx context.Context, req dto.ProductOptionCreateBulk) (api.Response, error) {
	if len(req.Options) == 0 {
		return nil, fmt.Errorf("no product options to create")
	}

	// Validate the product only once
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, fmt.Errorf("validation failed for product_id %s: %w", req.ProductId, err)
	}

	// Set productId into all options
	for i := range req.Options {
		req.Options[i].ProductId = req.ProductId
	}

	var created []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Insert(req.Options).
		Execute(&created); err != nil {
		return nil, fmt.Errorf("failed to create bulk product options: %v", err)
	}

	return api.Success(created), nil
}

func (s *productOptionService) UpdateBulkProductOption(ctx context.Context, req dto.ProductOptionBulkUpdate) (api.Response, error) {
	if len(req.Options) == 0 {
		return nil, fmt.Errorf("no product options to update")
	}

	// Validate the product exists
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	var updatedOptions []dto.ProductOption

	for _, option := range req.Options {
		// Ensure the option belongs to the correct productId
		updateData := map[string]interface{}{
			"name":       option.Name,
			"price":      option.Price,
			"detail":     option.Detail,
			"metadata":   option.Metadata,
			"updated_at": time.Now(),
		}

		var updated []dto.ProductOption
		if err := s.db.DB.
			From("product_options").
			Update(updateData).
			Eq("id", option.Id).
			Execute(&updated); err != nil {
			return nil, fmt.Errorf("failed to update product option ID %s: %v", option.Id, err)
		}

		updatedOptions = append(updatedOptions, updated...)
	}

	return api.Success(updatedOptions), nil
}
