package dto

import "time"

type Banner struct {
	Id        string    `json:"id"`
	Slot      string    `json:"slot"`
	BannerUrl string    `json:"banner_url"`
	Alt       string    `json:"alt"`
	Link      string    `json:"link"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerCreate struct {
	Slot      string `json:"slot" validate:"required,oneof=first second third"`
	BannerUrl string `json:"banner_url" validate:"required,url"`
	Alt       string `json:"alt"`
	Link      string `json:"link"`
	IsActive  bool   `json:"is_active"`
}

type BannerUpdate struct {
	Id        string `json:"id" validate:"required"`
	Slot      string `json:"slot" validate:"omitempty,oneof=first second third"`
	BannerUrl string `json:"banner_url" validate:"omitempty,url"`
	Alt       string `json:"alt"`
	Link      string `json:"link" validate:"omitempty,url"`
	IsActive  *bool  `json:"is_active"`
}

type BannerFilter struct {
	Slot     string `query:"slot"`
	IsActive *bool  `query:"is_active"`
}
