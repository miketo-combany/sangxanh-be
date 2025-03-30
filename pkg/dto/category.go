package dto

import (
	"SangXanh/pkg/enum"
	"github.com/google/uuid"
	"time"
)

type Category struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Metadata    map[string]interface{} `json:"metadata"`
	Status      bool                   `json:"status"`
	Thumbnail   string                 `json:"thumbnail"`
	Level       int                    `json:"level"`
	ParentId    string                 `json:"parent_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   time.Time              `json:"deleted_at"`
	Description string                 `json:"description"`
}

type CategoryCreate struct {
	Name        string                 `json:"name"`
	Thumbnail   string                 `json:"thumbnail"`
	ParentId    string                 `json:"parent_id,omitempty"`
	Status      bool                   `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	Description string                 `json:"description"`
	Level       int                    `json:"level"`
}

type CategoryUpdate struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Thumbnail   string                 `json:"thumbnail"`
	Status      bool                   `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	Description string                 `json:"description"`
}

type CategoryResponse struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Thumbnail   string                 `json:"thumbnail"`
	Level       int                    `json:"level"`
	Description string                 `json:"description"`
	Parent      *CategoryResponse      `json:"parent"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      enum.Status            `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CategoryListResponse struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Thumbnail   string                 `json:"thumbnail"`
	Level       int                    `json:"level"`
	Description string                 `json:"description"`
	Categories  []CategoryListResponse `json:"categories"`
	ParentId    string                 `json:"parent_id"`
	Status      enum.Status            `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func GetResponse(cate *Category, cateParent *Category) CategoryResponse {
	status := enum.Inactive
	if cate.Status {
		status = enum.Active
	}
	cateResponse := CategoryResponse{
		Id:          cate.Id,
		Name:        cate.Name,
		Thumbnail:   cate.Thumbnail,
		Level:       cate.Level,
		Description: cate.Description,
		UpdatedAt:   cate.UpdatedAt,
		CreatedAt:   cate.CreatedAt,
		Status:      status,
		Metadata:    cate.Metadata,
	}
	if cate.ParentId != uuid.New().String() && cateParent != nil {
		cateResponse.Parent = &CategoryResponse{
			Id:          cateParent.Id,
			Name:        cateParent.Name,
			Thumbnail:   cateParent.Thumbnail,
			Level:       cateParent.Level,
			Description: cateParent.Description,
		}
	}
	return cateResponse
}
