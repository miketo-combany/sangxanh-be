package dto

type Cart struct {
	ID              string `json:"id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	DeletedAt       string `json:"deleted_at"`
	UserID          string `json:"user_id"`
	ProductOptionID string `json:"product_option_id"`
	Quantity        int    `json:"quantity"`
}
