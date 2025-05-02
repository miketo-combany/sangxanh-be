package dto

import (
	"SangXanh/pkg/common/query"
	"SangXanh/pkg/enum"
	"time"
)

type Product struct {
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	Price        float32             `json:"price"`
	Content      string              `json:"content"`
	ImageDetail  string              `json:"image_detail"`
	Thumbnail    string              `json:"thumbnail"`
	CategoryId   string              `json:"category_id"`
	Discount     float32             `json:"discount"`
	DiscountType enum.DiscountType   `json:"discount_type"`
	Metadata     []map[string]string `json:"metadata"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    time.Time           `json:"deleted_at"`
}

type ProductCreated struct {
	Name         string              `json:"name"`
	Price        float32             `json:"price"`
	Content      string              `json:"content"`
	ImageDetail  string              `json:"image_detail"`
	Thumbnail    string              `json:"thumbnail"`
	CategoryId   string              `json:"category_id"`
	Discount     float32             `json:"discount"`
	DiscountType enum.DiscountType   `json:"discount_type"`
	Metadata     []map[string]string `json:"metadata"`
}

type ProductUpdated struct {
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	Price        float32             `json:"price"`
	Content      string              `json:"content"`
	ImageDetail  string              `json:"image_detail"`
	Thumbnail    string              `json:"thumbnail"`
	CategoryId   string              `json:"category_id"`
	Discount     float32             `json:"discount"`
	DiscountType enum.DiscountType   `json:"discount_type"`
	Metadata     []map[string]string `json:"metadata"`
}

type ProductResponse struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Price        float32           `json:"price"`
	Content      string            `json:"content"`
	ImageDetail  string            `json:"image_detail"`
	Thumbnail    string            `json:"thumbnail"`
	Discount     float32           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	CategoryId   string            `json:"category_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type ProductList struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Price        float64           `json:"price"`
	Content      string            `json:"content"`
	Thumbnail    string            `json:"thumbnail"`
	Category     CategoryProduct   `json:"categories"`
	Discount     float64           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	ImageDetail  string            `json:"image_detail"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type CategoryProduct struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProductFilter struct {
	query.Pagination
	Name        string  `query:"name"`
	CategoryId  string  `query:"category_id"`
	IsDiscount  bool    `query:"is_discount"`
	GreaterThan float64 `query:"greater_than"`
	SmallerThan float64 `query:"smaller_than"`
}

type ProductDetail struct {
	Id              string                  `json:"id"`
	Discount        float32                 `json:"discount"`
	Name            string                  `json:"name"`
	Price           float32                 `json:"price"`
	Content         string                  `json:"content"`
	ImageDetail     string                  `json:"image_detail"`
	Thumbnail       string                  `json:"thumbnail"`
	DiscountType    enum.DiscountType       `json:"discount_type"`
	MaxPrice        float32                 `json:"max_price"`
	MinPrice        float32                 `json:"min_price"`
	CategoryProduct CategoryProduct         `json:"categories"`
	ProductOptions  []ProductOptionResponse `json:"product_option_detail"`
	ProductVariants []ProductVariant        `json:"product_variant_detail"`
}
