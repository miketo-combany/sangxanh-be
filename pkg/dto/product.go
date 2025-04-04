package dto

import (
	"SangXanh/pkg/enum"
	"time"
)

type Product struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Price        float32           `json:"price"`
	Content      string            `json:"content"`
	ImageDetail  string            `json:"image_detail"`
	Thumbnail    string            `json:"thumbnail"`
	CategoryId   string            `json:"category_id"`
	Discount     float32           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    time.Time         `json:"deleted_at"`
}

type ProductCreated struct {
	Name         string            `json:"name"`
	Price        float32           `json:"price"`
	Content      string            `json:"content"`
	ImageDetail  string            `json:"image_detail"`
	Thumbnail    string            `json:"thumbnail"`
	CategoryId   string            `json:"category_id"`
	Discount     float32           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Metadata     map[string]string `json:"metadata"`
}

type ProductUpdated struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Price        float32           `json:"price"`
	Content      string            `json:"content"`
	ImageDetail  string            `json:"image_detail"`
	Thumbnail    string            `json:"thumbnail"`
	CategoryId   string            `json:"category_id"`
	Discount     float32           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Metadata     map[string]string `json:"metadata"`
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
}

type CategoryProduct struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProductFilter struct {
	CategoryId  string  `json:"category_id"`
	IsDiscount  bool    `json:"is_discount"`
	GreaterThan float64 `json:"greater_than"`
	SmallerThan float64 `json:"smaller_than"`
}
