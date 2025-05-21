package dto

import "SangXanh/pkg/common/query"

type Cart struct {
	ID              string `json:"id"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	DeletedAt       string `json:"deleted_at,omitempty"`
	UserID          string `json:"user_id"`
	ProductOptionID string `json:"product_option_id"`
	Quantity        int    `json:"quantity"`
}

type CartResponse struct {
	Cart
	ProductOption ProductOption `json:"product_option"`
}
type CartCreateRequest struct {
	UserID          string `json:"user_id"`
	ProductOptionID string `json:"product_option_id"`
	Quantity        int    `json:"quantity"`
}

type CartUpdate struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CartFilter struct {
	UserID string `json:"user_id"`
	query.Pagination
}
