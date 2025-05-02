package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"net/url"
	"strconv"
	"time"
)

type ProductService interface {
	ListProducts(ctx context.Context, filter dto.ProductFilter) (api.Response, error)
	CreateProduct(ctx context.Context, req dto.ProductCreated) (api.Response, error)
	UpdateProduct(ctx context.Context, req dto.ProductUpdated) (api.Response, error)
	DeleteProduct(ctx context.Context, id string) (api.Response, error)
	GetProductById(ctx context.Context, id string) (api.Response, error)
}

type productService struct {
	db *supabase.Client
}

func NewProductService(di do.Injector) (ProductService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UserService: %w", err)
	}

	return &productService{db: db}, nil
}

func (s *productService) countProducts(ctx context.Context, filter dto.ProductFilter) (int, error) {
	// Start the same base query you use in ListProducts
	q := s.db.DB.From("products").
		Select("id").
		IsNull("deleted_at") // keep soft‑deleted rows out
	// Apply the same filter conditions
	if filter.CategoryId != "" {
		q = q.Eq("category_id", filter.CategoryId)
	}
	if filter.IsDiscount {
		q = q.Not().IsNull("discount")
	}
	if filter.GreaterThan > 0 {
		q = q.Gt("price", strconv.FormatFloat(filter.GreaterThan, 'f', -1, 32))
	}
	if filter.SmallerThan > 0 {
		q = q.Lt("price", strconv.FormatFloat(filter.SmallerThan, 'f', -1, 32))
	}

	// Execute, discard the row data, keep the count
	var tmp []struct{} // dummy slice – we only care about the header that carries the count
	if err := q.Execute(&tmp); err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return len(tmp), nil
}

func (s *productService) ListProducts(ctx context.Context, filter dto.ProductFilter) (api.Response, error) {
	total, err := s.countProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 2. fetch the current page
	var products []dto.ProductList
	query := s.db.DB.From("products").
		Select("id,name,price,content,image_detail,category_id,thumbnail,discount,discount_type,categories!inner(id,name),created_at,updated_at").
		LimitWithOffset(int(filter.Limit), int((filter.Page-1)*filter.Limit)).
		IsNull("deleted_at")

	if filter.Name != "" {
		encoded := url.QueryEscape("%" + filter.Name + "%")
		query = query.Like("name", encoded)
	}
	if filter.CategoryId != "" {
		query = query.Eq("category_id", filter.CategoryId)
	}
	if filter.IsDiscount {
		query = query.Not().IsNull("discount")
	}
	if filter.GreaterThan > 0 {
		query = query.Gt("price", strconv.FormatFloat(filter.GreaterThan, 'f', -1, 32))
	}
	if filter.SmallerThan > 0 {
		query = query.Lt("price", strconv.FormatFloat(filter.SmallerThan, 'f', -1, 32))
	}

	if err := query.Execute(&products); err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	// 3. fill in pagination meta & return
	filter.Pagination.Total = int64(total)
	return api.SuccessPagination(products, &filter.Pagination), nil
}

func (s *productService) CreateProduct(ctx context.Context, req dto.ProductCreated) (api.Response, error) {
	newProduct := dto.ProductCreated{
		Name:         req.Name,
		Price:        req.Price,
		Content:      req.Content,
		ImageDetail:  req.ImageDetail,
		Thumbnail:    req.Thumbnail,
		CategoryId:   req.CategoryId,
		Discount:     req.Discount,
		DiscountType: req.DiscountType,
		Metadata:     req.Metadata,
	}

	err := s.validCategory(req.CategoryId)
	var product []dto.Product
	err = s.db.DB.From("products").Insert(newProduct).Execute(&product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	return api.Success(product[0]), nil
}

func (s *productService) UpdateProduct(ctx context.Context, req dto.ProductUpdated) (api.Response, error) {
	updateData := map[string]interface{}{
		"name":          req.Name,
		"price":         req.Price,
		"content":       req.Content,
		"image_detail":  req.ImageDetail,
		"thumbnail":     req.Thumbnail,
		"category_id":   req.CategoryId,
		"discount":      req.Discount,
		"discount_type": req.DiscountType,
		"metadata":      req.Metadata,
		"updated_at":    time.Now(),
	}
	err := s.validCategory(req.CategoryId)

	var product []dto.Product
	err = s.db.DB.From("products").Update(updateData).Eq("id", req.Id).Execute(&product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}
	return api.Success(product[0]), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) (api.Response, error) {
	updateData := map[string]interface{}{
		"deleted_at": time.Now(),
	}

	var product []dto.Product
	err := s.db.DB.From("products").Update(updateData).Eq("id", id).Execute(&product)
	if err != nil {
		return nil, fmt.Errorf("failed to soft delete product: %v", err)
	}
	return api.Success("Product deleted successfully"), nil
}

func (s *productService) validCategory(id string) error {
	var category []dto.Category
	err := s.db.DB.From("categories").Select("id").Eq("id", id).IsNull("deleted_at").Execute(&category)
	if err != nil {
		return fmt.Errorf("failed to find category: %v", err)
	}
	if len(category) == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// GetProductById returns a full product document (base info + category +
// option list + variant list).  All related rows must not be soft‑deleted.
func (s *productService) GetProductById(
	ctx context.Context,
	id string,
) (api.Response, error) {

	/*───────────────────────────────────────────────────────*
	 * 1)   Base product (+ category)                       *
	 *───────────────────────────────────────────────────────*/
	var rows []dto.ProductDetail
	if err := s.db.DB.
		From("products").
		Select("id,name,price,content,image_detail,category_id,thumbnail,discount,discount_type,categories!inner(id,name),created_at,updated_at").
		Eq("id", id).
		IsNull("deleted_at").
		Execute(&rows); err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("product not found")
	}
	product := rows[0]

	/*───────────────────────────────────────────────────────*
	 * 2)   Options (raw)                                   *
	 *───────────────────────────────────────────────────────*/
	var optRaw []dto.ProductOption
	if err := s.db.DB.
		From("product_options").
		Select("id,name,product_id,price,detail,metadata,created_at,updated_at").
		Eq("product_id", id).
		IsNull("deleted_at").
		Execute(&optRaw); err != nil {
		return nil, fmt.Errorf("failed to fetch product options: %w", err)
	}

	/*── gather variant_ids for a single name-lookup query ─*/
	var variantIDs []string
	for _, o := range optRaw {
		for _, d := range o.Detail {
			variantIDs = append(variantIDs, d.VariantId)
		}
	}
	nameByID, err := s.loadVariantNames(id, variantIDs)
	if err != nil {
		return nil, err
	}

	/*── enrich options + compute min/max price ────────────*/
	var (
		opts     []dto.ProductOptionResponse
		minPrice float64
		maxPrice float64
	)
	if len(optRaw) > 0 {
		minPrice, maxPrice = optRaw[0].Price, optRaw[0].Price
	}

	for _, o := range optRaw {
		if o.Price < minPrice {
			minPrice = o.Price
		}
		if o.Price > maxPrice {
			maxPrice = o.Price
		}

		view := dto.ProductOptionResponse{
			Id:        o.Id,
			Name:      o.Name,
			ProductId: o.ProductId,
			Price:     o.Price,
			Metadata:  o.Metadata,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
		}

		for _, d := range o.Detail {
			view.Detail = append(view.Detail, dto.ProductOptionVariantDetail{
				VariantId:    d.VariantId,
				VariantName:  nameByID[d.VariantId],
				VariantValue: d.Name, // same value you asked for
			})
		}
		opts = append(opts, view)
	}

	/*───────────────────────────────────────────────────────*
	 * 3)   Variants list (unchanged fetch)                 *
	 *───────────────────────────────────────────────────────*/
	var variants []dto.ProductVariant
	if err := s.db.DB.
		From("product_variants").
		Select("id,name,product_id,detail,metadata,created_at,updated_at").
		Eq("product_id", id).
		IsNull("deleted_at").
		Execute(&variants); err != nil {
		return nil, fmt.Errorf("failed to fetch product variants: %w", err)
	}

	/*───────────────────────────────────────────────────────*
	 * 4)   Assemble response                               *
	 *───────────────────────────────────────────────────────*/
	product.ProductOptions = opts
	product.ProductVariants = variants
	product.MinPrice = float32(minPrice)
	product.MaxPrice = float32(maxPrice)

	return api.Success(product), nil
}

// returns map[variantId]variantName (empty map if no IDs given)
func (s *productService) loadVariantNames(
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
