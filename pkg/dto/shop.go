package dto

import "time"

// Branch represents a shop branch structure
type Branch struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Hotline string `json:"hotline"`
}

// Policy represents a shop policy structure
type Policy struct {
	Link  string `json:"link"`
	Label string `json:"label"`
}

// Contact represents a contact information structure
type Contact struct {
	Link string `json:"link"`
	Type string `json:"type"`
}

// Shop represents the shop information response
type Shop struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Hotline   string    `json:"hotline"`
	Branches  []Branch  `json:"branches"`
	Policies  []Policy  `json:"policies"`
	Contact   []Contact `json:"contact"`
}

// ShopUpdate represents the request to update shop information
type ShopUpdate struct {
	Name     string    `json:"name" validate:"required"`
	Hotline  string    `json:"hotline" validate:"required"`
	Branches []Branch  `json:"branches" validate:"required,dive"`
	Policies []Policy  `json:"policies" validate:"required,dive"`
	Contact  []Contact `json:"contact" validate:"required,dive"`
}

// ShopResponse represents the response structure for shop operations
type ShopResponse struct {
	Shop Shop `json:"shop"`
}
