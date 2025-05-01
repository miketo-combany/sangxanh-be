package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/enum"
	"SangXanh/pkg/log"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"net/url"
	"time"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.CategoryCreate) (api.Response, error)
	ListCategories(ctx context.Context, req dto.ListCategory, name string) (api.Response, error)
	UpdateCategory(ctx context.Context, req dto.CategoryUpdate) (api.Response, error)
	DeleteCategory(ctx context.Context, categoryId string) (api.Response, error)
	ListCategoryById(ctx context.Context, categoryId string) (api.Response, error)
}

type categoryService struct {
	db *supabase.Client
}

func NewCategoryService(di do.Injector) (CategoryService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UserService: %w", err)
	}

	return &categoryService{db: db}, nil
}

func (u *categoryService) ListCategoryById(ctx context.Context, categoryId string) (api.Response, error) {
	// 1. Fetch the category that matches the given ID and has not been soft‑deleted
	var categories []dto.Category
	err := u.db.DB.
		From("categories").
		Select("*").
		Eq("id", categoryId).
		IsNull("deleted_at").
		Execute(&categories)

	if err != nil {
		log.Errorf("failed to fetch category %s: %v", categoryId, err)
		return nil, fmt.Errorf("failed to fetch category")
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("category not found")
	}
	cat := categories[0]
	var childCategories []dto.Category
	err = u.db.DB.From("categories").Select("*").Eq("parent_id", categoryId).Execute(&childCategories)

	// 2. Build the response payload (same fields you return elsewhere)
	categoryResponse := dto.CategoryResponse{
		Id:          cat.Id,
		Name:        cat.Name,
		Thumbnail:   cat.Thumbnail,
		Level:       cat.Level,
		Description: cat.Description,
		Status:      enum.ToStatus(cat.Status),
		Categories:  childCategories,
		Metadata:    cat.Metadata,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}

	return api.Success(categoryResponse), nil
}

func (u *categoryService) CreateCategory(ctx context.Context, req dto.CategoryCreate) (api.Response, error) {
	createCategory := dto.CategoryCreate{
		Name:              req.Name,
		Metadata:          req.Metadata,
		Status:            req.Status,
		Thumbnail:         req.Thumbnail,
		Description:       req.Description,
		IsDisplayHomepage: req.IsDisplayHomepage,
	}

	var parentCategory []dto.Category
	if req.ParentId != uuid.Nil.String() && req.ParentId != "" {
		err := u.db.DB.From("categories").
			Select("*").
			Eq("id", req.ParentId).
			Execute(&parentCategory)

		if err != nil {
			log.Errorf("Parent category with ID %s does not exist: %v", req.ParentId, err)
			return nil, err
		}

		createCategory.ParentId = req.ParentId
		createCategory.Level = parentCategory[0].Level + 1 // Set child category level
	} else {
		createCategory.ParentId = ""
	}

	var category []dto.Category
	err := u.db.DB.From("categories").Insert(createCategory).Execute(&category)
	if err != nil {
		log.Errorf("failed to insert category: %v", err)
		return nil, err
	}
	log.Info("category created", category)
	categoryResponse := dto.GetResponse(&category[0])
	return api.Success(categoryResponse), nil
}

func (u *categoryService) ListCategories(ctx context.Context, req dto.ListCategory, name string) (api.Response, error) {
	// Step 1: Fetch all categories from Supabase
	var categories []dto.Category
	query := u.db.DB.From("categories").Select("*").IsNull("deleted_at")
	if name != "" {
		encoded := url.QueryEscape("%" + name + "%")
		query = query.Like("name", encoded)
	}
	err := query.Execute(&categories)
	if err != nil {
		log.Errorf("failed to fetch categories: %v", err)
		return nil, err
	}

	// Step 2: Create a map to organize categories by ParentId
	categoryMap := make(map[string]dto.Category)
	categoryResponseMap := make(map[string]dto.CategoryListResponse)
	var topLevelCategories []dto.Category

	for _, category := range categories {
		if category.ParentId == uuid.Nil.String() || category.ParentId == "" {
			topLevelCategories = append(topLevelCategories, category) // Root categories
		}
		categoryMap[category.Id] = category // Group by ParentId
		categoryResponseMap[category.Id] = buildCategoryResponse(category)
	}

	// Step 3: Convert categories into the response format
	var categoryResponses []dto.CategoryListResponse

	for _, category := range categoryResponseMap {
		if category.ParentId != uuid.Nil.String() && category.ParentId != "" {
			parentCategory := categoryResponseMap[category.ParentId]
			parentCategory.Categories = append(parentCategory.Categories, category)
			categoryResponseMap[category.ParentId] = parentCategory
		}
	}

	for _, category := range topLevelCategories {
		categoryResponses = append(categoryResponses, categoryResponseMap[category.Id])
	}

	if int(req.Limit) == 0 && int(req.Page) == 0 {
		return api.Success(categoryResponses), nil
	}

	if int(req.Limit) >= len(categoryResponses) {
		return api.Success(categoryResponses), nil
	}

	if int(req.Limit*(req.Page-1)) > len(categoryResponses) {
		return api.SuccessPagination(nil, &req.Pagination), nil
	}

	req.Pagination.Total = int64(len(categoryResponses))

	categoryResponsesPage := categoryResponses[req.Limit*(req.Page-1) : req.Limit*req.Page]

	return api.SuccessPagination(categoryResponsesPage, &req.Pagination), nil
}

// Helper function to convert a single category
func buildCategoryResponse(category dto.Category) dto.CategoryListResponse {

	return dto.CategoryListResponse{
		Id:          category.Id,
		Name:        category.Name,
		Thumbnail:   category.Thumbnail,
		Level:       category.Level,
		Description: category.Description,
		ParentId:    category.ParentId,
		Status:      enum.ToStatus(category.Status),
		Metadata:    category.Metadata,
	}
}

func (u *categoryService) UpdateCategory(ctx context.Context, req dto.CategoryUpdate) (api.Response, error) {
	// Check if the category exists
	var existingCategory []dto.Category
	err := u.db.DB.From("categories").Select("*").Eq("id", req.Id).Execute(&existingCategory)
	if err != nil {
		log.Errorf("Category with ID %s not found: %v", req.Id, err)
		return nil, fmt.Errorf("category not found")
	}

	// Prepare updated fields
	updateData := map[string]interface{}{
		"name":                req.Name,
		"thumbnail":           req.Thumbnail,
		"status":              req.Status,
		"metadata":            req.Metadata,
		"description":         req.Description,
		"is_display_homepage": req.IsDisplayHomepage,
		"updated_at":          time.Now(),
	}

	// Perform the update
	var updateCategory []dto.Category
	err = u.db.DB.From("categories").Update(updateData).Eq("id", req.Id).Execute(&updateCategory)
	if err != nil {
		log.Errorf("Failed to update category %s: %v", req.Id, err)
		return nil, fmt.Errorf("failed to update category")
	}

	return api.Success(updateCategory[0]), nil
}

func (u *categoryService) DeleteCategory(ctx context.Context, categoryId string) (api.Response, error) {
	// Check if the category exists
	var category []dto.Category
	err := u.db.DB.From("categories").Select("*").Eq("id", categoryId).Execute(&category)
	if err != nil {
		log.Errorf("Category with ID %s not found: %v", categoryId, err)
		return nil, fmt.Errorf("category not found")
	}

	// Check if the category has child categories
	var childCategories []dto.Category
	err = u.db.DB.From("categories").Select("id").Eq("parent_id", categoryId).Execute(&childCategories)
	if err != nil {
		log.Errorf("Failed to check child categories for %s: %v", categoryId, err)
		return nil, fmt.Errorf("failed to verify child categories")
	}

	if len(childCategories) > 0 {
		log.Errorf("Cannot delete category %s as it has child categories", categoryId)
		return nil, fmt.Errorf("category has child categories and cannot be deleted")
	}

	updateData := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	// Delete the category
	err = u.db.DB.From("categories").Update(updateData).Eq("id", categoryId).Execute(&category)
	if err != nil {
		log.Errorf("Failed to delete category %s: %v", categoryId, err)
		return nil, fmt.Errorf("failed to delete category")
	}

	log.Infof("Category %s deleted successfully", categoryId)
	return api.Success("Category deleted successfully"), nil
}
