package dto

import (
	"SangXanh/pkg/common/query"
	"SangXanh/pkg/enum"
	"time"
)

type Category struct {
	Id                string                   `json:"id"`
	Name              string                   `json:"name"`
	Metadata          []map[string]interface{} `json:"metadata"`
	Status            bool                     `json:"status"`
	Thumbnail         string                   `json:"thumbnail"`
	Level             int                      `json:"level"`
	ParentId          string                   `json:"parent_id"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	DeletedAt         time.Time                `json:"deleted_at"`
	Description       string                   `json:"description"`
	IsDisplayHomepage bool                     `json:"is_display_homepage"`
	IsDisplayHeader   bool                     `json:"is_display_header"`
}

type CategoryCreate struct {
	Name              string                   `json:"name"`
	Thumbnail         string                   `json:"thumbnail"`
	ParentId          string                   `json:"parent_id,omitempty"`
	Status            bool                     `json:"status"`
	Metadata          []map[string]interface{} `json:"metadata"`
	Description       string                   `json:"description"`
	Level             int                      `json:"level"`
	IsDisplayHomepage bool                     `json:"is_display_homepage"`
	IsDisplayHeader   bool                     `json:"is_display_header"`
}

type CategoryUpdate struct {
	Id                string                   `json:"id"`
	Name              string                   `json:"name"`
	Thumbnail         string                   `json:"thumbnail"`
	Status            bool                     `json:"status"`
	Metadata          []map[string]interface{} `json:"metadata"`
	Description       string                   `json:"description"`
	IsDisplayHomepage bool                     `json:"is_display_homepage"`
	ParentId          string                   `json:"parent_id,omitempty"`
	IsDisplayHeader   bool                     `json:"is_display_header"`
}

type CategoryResponse struct {
	Id                string                   `json:"id"`
	Name              string                   `json:"name"`
	Thumbnail         string                   `json:"thumbnail"`
	Level             int                      `json:"level"`
	Description       string                   `json:"description"`
	Categories        []Category               `json:"categories"`
	UpdatedAt         time.Time                `json:"updated_at"`
	CreatedAt         time.Time                `json:"created_at"`
	Status            enum.Status              `json:"status"`
	Metadata          []map[string]interface{} `json:"metadata"`
	IsDisplayHomepage bool                     `json:"is_display_homepage"`
	IsDisplayHeader   bool                     `json:"is_display_header"`
}

type CategoryListResponse struct {
	Id                string                   `json:"id"`
	Name              string                   `json:"name"`
	Thumbnail         string                   `json:"thumbnail"`
	Level             int                      `json:"level"`
	Description       string                   `json:"description"`
	Categories        []CategoryListResponse   `json:"categories"`
	ParentId          string                   `json:"parent_id"`
	Status            enum.Status              `json:"status"`
	Metadata          []map[string]interface{} `json:"metadata"`
	IsDisplayHomepage bool                     `json:"is_display_homepage"`
	IsDisplayHeader   bool                     `json:"is_display_header"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

type ListCategory struct {
	query.Pagination
	IsDisplayHomepage bool   `query:"is_display_homepage"`
	IsDisplayHeader   bool   `query:"is_display_header"`
	Name              string `query:"name"`
}

func GetResponse(cate *Category) CategoryResponse {
	status := enum.Inactive
	if cate.Status {
		status = enum.Active
	}
	cateResponse := CategoryResponse{
		Id:                cate.Id,
		Name:              cate.Name,
		Thumbnail:         cate.Thumbnail,
		Level:             cate.Level,
		Description:       cate.Description,
		UpdatedAt:         cate.UpdatedAt,
		CreatedAt:         cate.CreatedAt,
		Status:            status,
		IsDisplayHomepage: cate.IsDisplayHomepage,
		IsDisplayHeader:   cate.IsDisplayHeader,
		Metadata:          cate.Metadata,
	}
	return cateResponse
}
