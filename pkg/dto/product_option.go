package dto

import "time"

type ProductOption struct {
	Id        string                `json:"id"`
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Metadata  []map[string]string   `json:"metadata"`
	Detail    []ProductOptionDetail `json:"detail"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type ProductOptionDetail struct {
	VariantId string `json:"variant_id"`
	Name      string `json:"name"`
}

type ProductOptionCreateBulk struct {
	ProductId string                `json:"product_id"`
	Options   []ProductOptionCreate `json:"options"`
}

type ProductOptionBulkUpdate struct {
	ProductId string                `json:"product_id"`
	Options   []ProductOptionUpdate `json:"options"`
}

type ProductOptionCreate struct {
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Detail    []ProductOptionDetail `json:"detail"`
	Metadata  []map[string]string   `json:"metadata"`
}

type ProductOptionUpdate struct {
	Id        string                `json:"id"`
	Name      string                `json:"name"`
	ProductId string                `json:"product_id"`
	Price     float64               `json:"price"`
	Detail    []ProductOptionDetail `json:"detail"`
	Metadata  []map[string]string   `json:"metadata"`
}

type ProductOptionResponse struct {
	Id        string                       `json:"id"`
	Name      string                       `json:"name"`
	ProductId string                       `json:"product_id"`
	Price     float64                      `json:"price"`
	Metadata  []map[string]string          `json:"metadata"`
	Detail    []ProductOptionVariantDetail `json:"detail"`
	CreatedAt time.Time                    `json:"created_at"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

type ProductOptionVariantDetail struct {
	VariantId    string `json:"variant_id"`
	VariantName  string `json:"variant_name"`
	VariantValue string `json:"variant_value"` // == Name from the request
}
