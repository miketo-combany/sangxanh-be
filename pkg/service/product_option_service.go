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

func (s *productOptionService) UpdateBulkProductOption(
	ctx context.Context,
	req dto.ProductOptionBulkUpdate,
) (api.Response, error) {

	// 0️⃣  Sanity checks -----------------------------------------------------
	if len(req.Options) == 0 {
		return nil, fmt.Errorf("no product options to update")
	}
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	// 1️⃣  Load current, non-deleted options for this product ---------------
	var current []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Select("id").
		Eq("product_id", req.ProductId).
		IsNull("deleted_at").
		Execute(&current); err != nil {
		return nil, fmt.Errorf("failed to fetch current product options: %v", err)
	}

	existingIDs := map[string]bool{}
	for _, opt := range current {
		existingIDs[opt.Id] = true
	}

	// 2️⃣  Iterate over the payload ----------------------------------------
	now := time.Now()
	payloadIDs := map[string]bool{}
	var result []dto.ProductOption

	for _, opt := range req.Options {
		if opt.Id == "" { // --- CREATE -------------------------------------
			opt.ProductId = req.ProductId
			createProduct := dto.ProductOptionCreate{
				ProductId: req.ProductId,
				Name:      opt.Name,
				Price:     opt.Price,
				Detail:    opt.Detail,
				Metadata:  opt.Metadata,
			}
			var created []dto.ProductOption
			if err := s.db.DB.
				From("product_options").
				Insert(createProduct).
				Execute(&created); err != nil {
				return nil, fmt.Errorf("failed to create option %q: %v", opt.Name, err)
			}
			result = append(result, created...)
			continue
		}

		// --- UPDATE --------------------------------------------------------
		payloadIDs[opt.Id] = true

		updateData := map[string]interface{}{
			"name":       opt.Name,
			"price":      opt.Price,
			"detail":     opt.Detail,
			"metadata":   opt.Metadata,
			"updated_at": now,
		}

		var updated []dto.ProductOption
		if err := s.db.DB.
			From("product_options").
			Update(updateData).
			Eq("id", opt.Id).
			Eq("product_id", req.ProductId). // extra guard
			Execute(&updated); err != nil {
			return nil, fmt.Errorf("failed to update option %s: %v", opt.Id, err)
		}
		result = append(result, updated...)
	}

	// 3️⃣  Soft-delete rows not included in the payload ---------------------
	for id := range existingIDs {
		if !payloadIDs[id] {
			if err := s.db.DB.
				From("product_options").
				Update(map[string]interface{}{"deleted_at": now}).
				Eq("id", id).
				Execute(nil); err != nil { // no RETURNING needed
				return nil, fmt.Errorf("failed to soft-delete option %s: %v", id, err)
			}
		}
	}

	return api.Success(result), nil
}
