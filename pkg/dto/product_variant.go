package dto

import "time"

type ProductVariant struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	ProductId string                 `json:"product_id"`
	Detail    []ProductVariantDetail `json:"detail"`
	Metadata  map[string]string      `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type ProductVariantDetail struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductVariantResponse struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	Product   Product                `json:"products"`
	Detail    []ProductVariantDetail `json:"detail"`
	Metadata  map[string]string      `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type ProductVariantCreate struct {
	Name      string `json:"name"`
	ProductId string `json:"product_id"`
	Detail    string `json:"detail"`
	Metadata  string `json:"metadata"`
}

type ProductVariantUpdate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Detail   string `json:"detail"`
	Metadata string `json:"metadata"`
}
