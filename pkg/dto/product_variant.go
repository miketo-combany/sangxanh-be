package dto

import "time"

type ProductVariant struct {
	Id        string                   `json:"id"`
	Name      string                   `json:"name"`
	ProductId string                   `json:"product_id"`
	Detail    []ProductVariantDetail   `json:"detail"`
	Metadata  []map[string]interface{} `json:"metadata"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type ProductVariantDetail struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductVariantResponse struct {
	Id        string                   `json:"id"`
	Name      string                   `json:"name"`
	Product   Product                  `json:"products"`
	Detail    []ProductVariantDetail   `json:"detail"`
	Metadata  []map[string]interface{} `json:"metadata"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type ProductVariantCreate struct {
	Name      string                   `json:"name"`
	ProductId string                   `json:"product_id"`
	Detail    []ProductVariantDetail   `json:"detail"`
	Metadata  []map[string]interface{} `json:"metadata"`
}

type ProductVariantCreateBulk struct {
	ProductId string                 `json:"product_id"`
	Variants  []ProductVariantCreate `json:"variants"`
}

type ProductVariantUpdateBulk struct {
	ProductId string                 `json:"product_id"`
	Variants  []ProductVariantUpdate `json:"variants"`
}

type ProductVariantUpdate struct {
	Id       string                   `json:"id"`
	Name     string                   `json:"name"`
	Detail   []ProductVariantDetail   `json:"detail"`
	Metadata []map[string]interface{} `json:"metadata"`
}
