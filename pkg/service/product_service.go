package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"strconv"
	"time"
)

type ProductService interface {
	ListProducts(ctx context.Context, filter dto.ProductFilter) (api.Response, error)
	CreateProduct(ctx context.Context, req dto.ProductCreated) (api.Response, error)
	UpdateProduct(ctx context.Context, req dto.ProductUpdated) (api.Response, error)
	DeleteProduct(ctx context.Context, id string) (api.Response, error)
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

func (s *productService) ListProducts(ctx context.Context, filter dto.ProductFilter) (api.Response, error) {
	var products []dto.ProductList
	query := s.db.DB.From("products").
		Select("id,name,price,content,image_detail,category_id,thumbnail,discount,discount_type,categories!inner(id,name)").
		IsNull("deleted_at")

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

	err := query.Execute(&products)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}

	return api.Success(products), nil
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

	return api.Success(newProduct), nil
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
	}
	err := s.validCategory(req.CategoryId)

	var product []dto.Product
	err = s.db.DB.From("products").Update(updateData).Eq("id", req.Id).Execute(product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}
	return api.Success("Product updated successfully"), nil
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
