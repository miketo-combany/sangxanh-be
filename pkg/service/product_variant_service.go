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
// If id == "" (or == "ALL") we mass-delete every variant for product_id.
// DeleteProductVariant receives ONE variant ID, figures out which product it
// belongs to, and then soft-deletes *all* variants + *all* options for that
// product in a single transaction-style operation.
func (s *productVariantService) DeleteProductVariant(
	ctx context.Context, id string,
) (api.Response, error) {

	// ── 0. Guard clause ────────────────────────────────────────────────────
	if id == "" {
		return nil, fmt.Errorf("variant id must not be empty")
	}

	// ── 1. Look up the product_id for this variant ────────────────────────
	var row []struct {
		ProductId string `json:"product_id"`
	}
	if err := s.db.DB.
		From("product_variants").
		Select("product_id").
		Eq("id", id).
		Execute(&row); err != nil {
		return nil, fmt.Errorf("variant not found: %v", err)
	}
	productID := row[0].ProductId

	// ── 2. Timestamp to apply to all records we are marking deleted ───────
	now := time.Now()

	// ── 3. Soft-delete ALL variants for that product ──────────────────────
	if err := s.db.DB.
		From("product_variants").
		Update(map[string]interface{}{"deleted_at": now}).
		Eq("id", id).
		IsNull("deleted_at").
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to delete product variants: %v", err)
	}

	// ── 4. Soft-delete ALL product options for that product ───────────────
	if err := s.db.DB.
		From("product_options").
		Update(map[string]interface{}{"deleted_at": now}).
		Eq("product_id", productID).
		IsNull("deleted_at").
		Execute(nil); err != nil {
		return nil, fmt.Errorf("failed to delete product options: %v", err)
	}

	return api.Success("All variants and options for the product were deleted successfully"), nil
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

// UpdateBulkProductVariant creates / updates the rows received in req.Variants
// and soft-deletes any missing ones.  If at least one variant is deleted,
// every product_option for the same product_id is also soft-deleted.
func (s *productVariantService) UpdateBulkProductVariant(
	ctx context.Context,
	req dto.ProductVariantUpdateBulk,
) (api.Response, error) {

	// 0️⃣  Basic validation --------------------------------------------------
	if len(req.Variants) == 0 {
		return nil, fmt.Errorf("no product variants to update")
	}
	if err := s.validProduct(req.ProductId); err != nil {
		return nil, err
	}

	// 1️⃣  Fetch current, non-deleted variant IDs ---------------------------
	var current []dto.ProductVariant
	if err := s.db.DB.
		From("product_variants").
		Select("id").
		Eq("product_id", req.ProductId).
		IsNull("deleted_at").
		Execute(&current); err != nil {
		return nil, fmt.Errorf("failed to fetch current product variants: %v", err)
	}

	existingIDs := make(map[string]struct{}, len(current))
	for _, v := range current {
		existingIDs[v.Id] = struct{}{}
	}

	// 2️⃣  Walk through the payload ----------------------------------------
	now := time.Now()
	payloadIDs := map[string]struct{}{}
	var result []dto.ProductVariant

	for _, v := range req.Variants {

		/* ---------- CREATE ------------------------------------------------ */
		if v.Id == "" {
			createReq := dto.ProductVariantCreate{
				Name:      v.Name,
				ProductId: req.ProductId,
				Detail:    v.Detail,
				Metadata:  v.Metadata,
			}

			var created []dto.ProductVariant
			if err := s.db.DB.
				From("product_variants").
				Insert(createReq).
				Execute(&created); err != nil {
				return nil, fmt.Errorf("failed to create variant %q: %v", v.Name, err)
			}
			result = append(result, created...)
			continue
		}

		/* ---------- UPDATE ------------------------------------------------ */
		payloadIDs[v.Id] = struct{}{}

		updateData := map[string]interface{}{
			"name":       v.Name,
			"detail":     v.Detail,
			"metadata":   v.Metadata,
			"updated_at": now,
		}

		var updated []dto.ProductVariant
		if err := s.db.DB.
			From("product_variants").
			Update(updateData).
			Eq("id", v.Id).
			Eq("product_id", req.ProductId). // safety guard
			Execute(&updated); err != nil {
			return nil, fmt.Errorf("failed to update variant %s: %v", v.Id, err)
		}
		result = append(result, updated...)
	}

	// 3️⃣  Soft-delete rows missing from the payload -----------------------
	needOptionCleanup := false

	for id := range existingIDs {
		if _, keep := payloadIDs[id]; !keep {
			if err := s.db.DB.
				From("product_variants").
				Update(map[string]interface{}{"deleted_at": now}).
				Eq("id", id).
				Execute(nil); err != nil {
				return nil, fmt.Errorf("failed to soft-delete variant %s: %v", id, err)
			}
			needOptionCleanup = true
		}
	}

	// 4️⃣  Cascade-style cleanup of product_options ------------------------
	if needOptionCleanup {
		if err := s.db.DB.
			From("product_options").
			Update(map[string]interface{}{"deleted_at": now}).
			Eq("product_id", req.ProductId).
			IsNull("deleted_at").
			Execute(nil); err != nil {
			return nil, fmt.Errorf("failed to soft-delete product options for product %s: %v", req.ProductId, err)
		}
	}

	return api.Success(result), nil
}

// deleteAllVariantsByProduct sets deleted_at for every non-deleted row
// belonging to the given product.  It is idempotent.
func (s *productVariantService) deleteAllVariantsByProduct(
	ctx context.Context,
	productID string,
) error {
	err := s.db.DB.
		From("product_variants").
		Update(map[string]interface{}{"deleted_at": time.Now()}).
		Eq("product_id", productID).
		IsNull("deleted_at").
		Execute(nil) // we don’t need RETURNING
	return err
}
