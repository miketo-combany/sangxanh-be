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

type ProductVariantService interface {
	ListProductVariants(ctx context.Context, productId string) (api.Response, error)
	CreateProductVariant(ctx context.Context, req dto.ProductVariantCreate) (api.Response, error)
	UpdateProductVariant(ctx context.Context, req dto.ProductVariantUpdate) (api.Response, error)
	DeleteProductVariant(ctx context.Context, id string) (api.Response, error)
	CreateBulkProductVariant(ctx context.Context, reqs dto.ProductVariantCreateBulk) (api.Response, error)
	UpdateBulkProductVariant(ctx context.Context, reqs dto.ProductVariantUpdateBulk) (api.Response, error)
}

type productVariantService struct {
	db *supabase.Client
}

func NewProductVariantService(di do.Injector) (ProductVariantService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ProductVariantService: %w", err)
	}

	return &productVariantService{db: db}, nil
}

// ListProductVariants fetches all variants that have not been soft-deleted (deleted_at IS NULL).
func (s *productVariantService) ListProductVariants(ctx context.Context, productId string) (api.Response, error) {
	var variants []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Select("id,name,product_id,detail,metadata,created_at,updated_at").
		Eq("product_id", productId).
		IsNull("deleted_at").
		Execute(&variants)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product variants: %v", err)
	}

	return api.Success(variants), nil
}

// CreateProductVariant inserts a new product variant into the database.
func (s *productVariantService) CreateProductVariant(ctx context.Context, req dto.ProductVariantCreate) (api.Response, error) {
	// Ensure the referenced product exists (similar to validCategory in your product service).
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	var created []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Insert(req).
		Execute(&created)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %v", err)
	}

	return api.Success(created[0]), nil
}

// UpdateProductVariant updates an existing product variant's fields.
// It unmarshals the `detail` and `metadata` JSON before updating.
func (s *productVariantService) UpdateProductVariant(ctx context.Context, req dto.ProductVariantUpdate) (api.Response, error) {
	updateData := map[string]interface{}{
		"name":       req.Name,
		"detail":     req.Detail,
		"metadata":   req.Metadata,
		"updated_at": time.Now(),
	}

	var updated []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Update(updateData).
		Eq("id", req.Id).
		Execute(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update product variant: %v", err)
	}

	return api.Success(updated[0]), nil
}

// DeleteProductVariant performs a soft delete by setting deleted_at,
// similar to your DeleteProduct method for products.
func (s *productVariantService) DeleteProductVariant(ctx context.Context, id string) (api.Response, error) {
	updateData := map[string]interface{}{
		"deleted_at": time.Now(),
	}

	var updated []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Update(updateData).
		Eq("id", id).
		Execute(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product variant: %v", err)
	}

	return api.Success("Product variant deleted successfully"), nil
}

// validProduct checks that the product_id provided in a variant actually exists
// and has not been soft-deleted.
func (s *productVariantService) validProduct(id string) error {
	var product []dto.Product
	err := s.db.DB.From("products").
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

func (s *productVariantService) CreateBulkProductVariant(ctx context.Context, reqs dto.ProductVariantCreateBulk) (api.Response, error) {
	// Validate if product exists
	if err := s.validProduct(reqs.ProductId); err != nil {
		return nil, err
	}

	// Attach productId to all variants (safety)
	for i := range reqs.Variants {
		reqs.Variants[i].ProductId = reqs.ProductId
	}

	var created []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Insert(reqs.Variants).
		Execute(&created)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk create product variants: %v", err)
	}

	return api.Success(created), nil
}

func (s *productVariantService) UpdateBulkProductVariant(ctx context.Context, reqs dto.ProductVariantUpdateBulk) (api.Response, error) {
	// Validate if product exists
	if err := s.validProduct(reqs.ProductId); err != nil {
		return nil, err
	}

	updatedVariants := []dto.ProductVariant{}

	for _, req := range reqs.Variants {
		updateData := map[string]interface{}{
			"name":       req.Name,
			"detail":     req.Detail,
			"metadata":   req.Metadata,
			"updated_at": time.Now(),
		}

		var updated []dto.ProductVariant
		err := s.db.DB.From("product_variants").
			Update(updateData).
			Eq("id", req.Id).
			Eq("product_id", reqs.ProductId). // make sure belongs to correct product
			Execute(&updated)
		if err != nil {
			return nil, fmt.Errorf("failed to update product variant with id %s: %v", req.Id, err)
		}

		updatedVariants = append(updatedVariants, updated[0])
	}

	return api.Success(updatedVariants), nil
}
