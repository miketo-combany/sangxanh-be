package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
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

func NewProductService(db *supabase.Client) ProductService {
	return &productService{db: db}
}

func (s *productService) ListProducts(ctx context.Context, filter dto.ProductFilter) (api.Response, error) {
	var products []dto.ProductList
	query := s.db.DB.From("products").Select("products.*, categories.name as category_name").
		Wfts("categories", "categories.id = products.category_id").
		IsNull("products.deleted_at")

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

	var product []dto.Product
	err := s.db.DB.From("products").Insert(newProduct).Execute(&product)
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

	err := s.db.DB.From("products").Update(updateData).Eq("id", req.Id).Execute()
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
