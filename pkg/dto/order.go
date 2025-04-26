package dto

import (
	"SangXanh/pkg/enum"
	"time"
)

type Order struct {
	Id        string                   `json:"id"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
	DeletedAt time.Time                `json:"deleted_at"`
	UserId    string                   `json:"user_id"`
	Address   string                   `json:"address"`
	Status    enum.OrderStatus         `json:"status"`
	Metadata  []map[string]interface{} `json:"metadata"`
}

type OrderDetail struct {
	Id              string                   `json:"id"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       time.Time                `json:"deleted_at"`
	OrderId         string                   `json:"order_id"`
	ProductOptionId string                   `json:"product_option_id"`
	Quantity        int                      `json:"quantity"`
	Discount        float64                  `json:"discount"`
	DiscountType    enum.DiscountType        `json:"discount_type"`
	Metadata        []map[string]interface{} `json:"metadata"`
}

type OrderDetailBase struct {
	Id              string                   `json:"id"`
	ProductOptionId string                   `json:"product_option_id"`
	Quantity        int                      `json:"quantity"`
	Discount        float64                  `json:"discount"`
	DiscountType    enum.DiscountType        `json:"discount_type"`
	Metadata        []map[string]interface{} `json:"metadata"`
}

type OrderCreate struct {
	UserId   string                   `json:"user_id"`
	Address  string                   `json:"address"`
	Status   enum.OrderStatus         `json:"status"`
	Metadata []map[string]interface{} `json:"metadata"`
}
type OrderUpdate struct {
	Id           string                   `json:"id"`
	UserId       string                   `json:"user_id"`
	Address      string                   `json:"address"`
	Status       enum.OrderStatus         `json:"status"`
	Metadata     []map[string]interface{} `json:"metadata"`
	OrderDetails []OrderDetailBase        `json:"order_details"`
}

type OrderDetailResponse struct {
	Order
	OrderDetail []OrderDetail `json:"order_detail"`
}

type OrderListFilter struct {
	Status enum.OrderStatus `json:"status"`
}
