package dto

import "time"

type ProductOption struct {
	Id        string                `json:"id"`
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Metadata  map[string]string     `json:"metadata"`
	Detail    []ProductOptionDetail `json:"detail"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type ProductOptionDetail struct {
	VariantId string `json:"variant_id"`
	Name      string `json:"name"`
}

type ProductOptionCreate struct {
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Detail    []ProductOptionDetail `json:"detail"`
	Metadata  map[string]string     `json:"metadata"`
}

type ProductOptionUpdate struct {
	Id        string                `json:"id"`
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Detail    []ProductOptionDetail `json:"detail"`
	Metadata  map[string]string     `json:"metadata"`
}
