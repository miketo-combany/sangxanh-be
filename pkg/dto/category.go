package dto

import (
	"SangXanh/pkg/common/query"
	"SangXanh/pkg/enum"
	"time"
)

type Category struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Icon              string                 `json:"icon"`
	Metadata          map[string]interface{} `json:"metadata"`
	Status            bool                   `json:"status"`
	Thumbnail         string                 `json:"thumbnail"`
	Level             int                    `json:"level"`
	ParentId          string                 `json:"parent_id"`
	FavoProductIds    []string               `json:"favo_product_ids"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         time.Time              `json:"deleted_at"`
	Description       string                 `json:"description"`
	IsDisplayHomepage bool                   `json:"is_display_homepage"`
	IsDisplayHeader   bool                   `json:"is_display_header"`
	IsSubHeader       bool                   `json:"is_sub_header"`
}

type ProductShortInfo struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Thumbnail  string  `json:"thumbnail"`
	Price      float64 `json:"price"`
	CategoryId string  `json:"category_id"`
}

type CategoryCreate struct {
	Name              string                 `json:"name"`
	Icon              string                 `json:"icon"`
	Thumbnail         string                 `json:"thumbnail"`
	ParentId          string                 `json:"parent_id,omitempty"`
	Status            bool                   `json:"status"`
	Metadata          map[string]interface{} `json:"metadata"`
	Description       string                 `json:"description"`
	Level             int                    `json:"level"`
	IsDisplayHomepage bool                   `json:"is_display_homepage"`
	IsDisplayHeader   bool                   `json:"is_display_header"`
	IsSubHeader       bool                   `json:"is_sub_header"`
	FavoProductIds    []string               `json:"favo_product_ids"`
}

type CategoryUpdate struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Icon              string                 `json:"icon"`
	Thumbnail         string                 `json:"thumbnail"`
	Status            bool                   `json:"status"`
	Metadata          map[string]interface{} `json:"metadata"`
	Description       string                 `json:"description"`
	IsDisplayHomepage bool                   `json:"is_display_homepage"`
	ParentId          string                 `json:"parent_id,omitempty"`
	IsDisplayHeader   bool                   `json:"is_display_header"`
	IsSubHeader       bool                   `json:"is_sub_header"`
	FavoProductIds    []string               `json:"favo_product_ids"`
}

type CategoryResponse struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Thumbnail         string                 `json:"thumbnail"`
	Level             int                    `json:"level"`
	Icon              string                 `json:"icon"`
	Description       string                 `json:"description"`
	Categories        []Category             `json:"categories"`
	FavoProductIds    []string               `json:"favo_product_ids"`
	FavoProducts      []ProductShortInfo     `json:"favo_products"`
	UpdatedAt         time.Time              `json:"updated_at"`
	CreatedAt         time.Time              `json:"created_at"`
	Status            enum.Status            `json:"status"`
	Metadata          map[string]interface{} `json:"metadata"`
	IsDisplayHomepage bool                   `json:"is_display_homepage"`
	IsDisplayHeader   bool                   `json:"is_display_header"`
	IsSubHeader       bool                   `json:"is_sub_header"`
}

type CategoryListResponse struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Icon              string                 `json:"icon"`
	Thumbnail         string                 `json:"thumbnail"`
	Level             int                    `json:"level"`
	Description       string                 `json:"description"`
	Categories        []CategoryListResponse `json:"categories"`
	ParentId          string                 `json:"parent_id"`
	Status            enum.Status            `json:"status"`
	Metadata          map[string]interface{} `json:"metadata"`
	IsDisplayHomepage bool                   `json:"is_display_homepage"`
	IsDisplayHeader   bool                   `json:"is_display_header"`
	IsSubHeader       bool                   `json:"is_sub_header"`
	FavoProductIds    []string               `json:"favo_product_ids"`
	FavoProducts      []ProductShortInfo     `json:"favo_products"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type ListCategory struct {
	query.Pagination
	IsDisplayHomepage bool   `query:"is_display_homepage"`
	IsDisplayHeader   bool   `query:"is_display_header"`
	IsSubHeader       bool   `query:"is_sub_header"`
	Name              string `query:"name"`
	ParentId          string `query:"parent_id"`
}

func GetResponse(cate *Category) CategoryResponse {
	status := enum.Inactive
	if cate.Status {
		status = enum.Active
	}
	cateResponse := CategoryResponse{
		Id:                cate.Id,
		Icon:              cate.Icon,
		Name:              cate.Name,
		Thumbnail:         cate.Thumbnail,
		Level:             cate.Level,
		Description:       cate.Description,
		UpdatedAt:         cate.UpdatedAt,
		CreatedAt:         cate.CreatedAt,
		Status:            status,
		IsDisplayHomepage: cate.IsDisplayHomepage,
		IsDisplayHeader:   cate.IsDisplayHeader,
		IsSubHeader:       cate.IsSubHeader,
		Metadata:          cate.Metadata,
		FavoProductIds:    cate.FavoProductIds,
	}
	return cateResponse
}
