package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"context"
	"fmt"

	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

type BannerService interface {
	ListBanners(ctx context.Context, filter dto.BannerFilter) (api.Response, error)
	CreateBanner(ctx context.Context, req dto.BannerCreate) (api.Response, error)
	UpdateBanner(ctx context.Context, req dto.BannerUpdate) (api.Response, error)
	DeleteBanner(ctx context.Context, id string) (api.Response, error)
	GetBannerById(ctx context.Context, id string) (api.Response, error)
}

type bannerService struct {
	db *supabase.Client
}

func NewBannerService(di do.Injector) (BannerService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize BannerService: %w", err)
	}
	return &bannerService{db: db}, nil
}

func (s *bannerService) ListBanners(ctx context.Context, filter dto.BannerFilter) (api.Response, error) {
	var banners []dto.Banner

	q := s.db.DB.From("banners").Select("*")

	// Apply filters
	if filter.Slot != "" {
		q.Eq("slot", filter.Slot)
	}

	if filter.IsActive != nil {
		if *filter.IsActive {
			q.Eq("is_active", "true")
		} else {
			q.Eq("is_active", "false")
		}
	}

	if err := q.OrderBy("created_at", "desc").Execute(&banners); err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	return api.Success(banners), nil
}

func (s *bannerService) CreateBanner(ctx context.Context, req dto.BannerCreate) (api.Response, error) {
	var result []dto.Banner

	// Set default value for is_active if not provided
	if !req.IsActive {
		req.IsActive = true
	}

	err := s.db.DB.From("banners").
		Insert(req).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no banner was created")
	}

	return api.Success(result[0]), nil
}

func (s *bannerService) UpdateBanner(ctx context.Context, req dto.BannerUpdate) (api.Response, error) {
	var result []dto.Banner

	updateData := make(map[string]interface{})

	if req.Slot != "" {
		updateData["slot"] = req.Slot
	}
	if req.BannerUrl != "" {
		updateData["banner_url"] = req.BannerUrl
	}
	if req.Alt != "" {
		updateData["alt"] = req.Alt
	}
	if req.Link != "" {
		updateData["link"] = req.Link
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}

	if len(updateData) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	err := s.db.DB.From("banners").
		Update(updateData).
		Eq("id", req.Id).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("banner not found")
	}

	return api.Success(result[0]), nil
}

func (s *bannerService) DeleteBanner(ctx context.Context, id string) (api.Response, error) {
	var result []dto.Banner

	err := s.db.DB.From("banners").
		Delete().
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to delete banner: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("banner not found")
	}

	return api.Success(map[string]string{"message": "Banner deleted successfully"}), nil
}

func (s *bannerService) GetBannerById(ctx context.Context, id string) (api.Response, error) {
	var result []dto.Banner

	err := s.db.DB.From("banners").
		Select("*").
		Eq("id", id).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("banner not found")
	}

	return api.Success(result[0]), nil
}
