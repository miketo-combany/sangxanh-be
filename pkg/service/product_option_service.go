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

func (s *productOptionService) ListProductOptions(
	ctx context.Context,
	productId string,
) (api.Response, error) {

	// ─── ① fetch options (same as before) ──────────────────────────────────
	var raw []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Select("id,name,product_id,price,detail,metadata,created_at,updated_at").
		Eq("product_id", productId).
		IsNull("deleted_at").
		Execute(&raw); err != nil {
		return nil, fmt.Errorf("failed to fetch product options: %w", err)
	}

	// ─── ② collect every variant_id we saw ─────────────────────────────────
	var variantIDs []string
	for _, opt := range raw {
		for _, d := range opt.Detail {
			variantIDs = append(variantIDs, d.VariantId)
		}
	}

	// ─── ③ one query to resolve names ─────────────────────────────────────
	nameByID, err := s.loadVariantNames(productId, variantIDs)
	if err != nil {
		return nil, err
	}

	// ─── ④ build response DTOs ────────────────────────────────────────────
	out := make([]dto.ProductOptionResponse, 0, len(raw))
	for _, opt := range raw {
		rsp := dto.ProductOptionResponse{
			Id:        opt.Id,
			Name:      opt.Name,
			ProductId: opt.ProductId,
			Price:     opt.Price,
			Metadata:  opt.Metadata,
			CreatedAt: opt.CreatedAt,
			UpdatedAt: opt.UpdatedAt,
		}
		for _, d := range opt.Detail {
			rsp.Detail = append(rsp.Detail, dto.ProductOptionVariantDetail{
				VariantId:    d.VariantId,
				VariantName:  nameByID[d.VariantId], // “Color”, “Size”,…
				VariantValue: d.Name,                // payload’s Name field
			})
		}
		out = append(out, rsp)
	}

	return api.Success(out), nil
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

	var allIds []string
	for _, opt := range req.Options {
		for _, d := range opt.Detail {
			allIds = append(allIds, d.VariantId)
		}
	}
	_, err := s.fetchAndValidateVariants(req.ProductId, allIds)
	if err != nil {
		return nil, err
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
	var allIds []string
	for _, opt := range req.Options {
		for _, d := range opt.Detail {
			allIds = append(allIds, d.VariantId)
		}
	}
	_, err := s.fetchAndValidateVariants(req.ProductId, allIds)
	if err != nil {
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

func (s *productOptionService) fetchAndValidateVariants(
	productId string,
	variantIDs []string,
) (map[string]string, error) {

	// de-duplicate to keep the SQL IN (…) short
	uniq := make(map[string]struct{}, len(variantIDs))
	for _, id := range variantIDs {
		uniq[id] = struct{}{}
	}
	ids := make([]string, 0, len(uniq))
	for id := range uniq {
		ids = append(ids, id)
	}

	// ---- single query ----------------------------------------------------
	var rows []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	if err := s.db.DB.
		From("product_variants").
		Select("id,name").
		Eq("product_id", productId).
		In("id", ids).
		IsNull("deleted_at").
		Execute(&rows); err != nil {
		return nil, fmt.Errorf("failed to load variants: %w", err)
	}

	names := make(map[string]string, len(rows))
	for _, r := range rows {
		names[r.Id] = r.Name
		delete(uniq, r.Id) // remove anything we actually found
	}

	// anything still in uniq is missing in DB ------------------------------
	if len(uniq) > 0 {
		missing := make([]string, 0, len(uniq))
		for id := range uniq {
			missing = append(missing, id)
		}
		return names, fmt.Errorf("unknown variant_id(s): %v", missing)
	}

	return names, nil // every id is legit
}

// returns map[variantId]variantName (empty map if no IDs given)
func (s *productOptionService) loadVariantNames(
	productId string,
	variantIDs []string,
) (map[string]string, error) {

	if len(variantIDs) == 0 {
		return map[string]string{}, nil
	}

	// deduplicate so the SQL IN (…) stays short
	uniq := make(map[string]struct{}, len(variantIDs))
	for _, id := range variantIDs {
		uniq[id] = struct{}{}
	}
	ids := make([]string, 0, len(uniq))
	for id := range uniq {
		ids = append(ids, id)
	}

	var rows []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	if err := s.db.DB.
		From("product_variants").
		Select("id,name").
		Eq("product_id", productId).
		In("id", ids).
		IsNull("deleted_at").
		Execute(&rows); err != nil {
		return nil, fmt.Errorf("failed to load variant names: %w", err)
	}

	out := make(map[string]string, len(rows))
	for _, r := range rows {
		out[r.Id] = r.Name
	}
	return out, nil
}
