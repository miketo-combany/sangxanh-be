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
// It unmarshals the `detail` and `metadata` JSON strings into their respective fields.
func (s *productVariantService) CreateProductVariant(ctx context.Context, req dto.ProductVariantCreate) (api.Response, error) {
	// Ensure the referenced product exists (similar to validCategory in your product service).
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	//var metadata map[string]string
	//if err := json.Unmarshal([]byte(req.Metadata), &metadata); err != nil {
	//	return nil, fmt.Errorf("invalid metadata format: %v", err)
	//}

	var created []dto.ProductVariant
	err := s.db.DB.From("product_variants").
		Insert(req).
		Execute(&created)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %v", err)
	}

	return api.Success(created), nil
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

	return api.Success("Product variant updated successfully"), nil
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
