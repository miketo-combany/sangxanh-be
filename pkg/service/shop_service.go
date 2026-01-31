package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/database"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

type ShopService interface {
	GetShop(ctx context.Context) (api.Response, error)
	UpdateShop(ctx context.Context, req dto.ShopUpdate) (api.Response, error)
}

type shopService struct {
	db *supabase.Client
}

func NewShopService(di do.Injector) (ShopService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ShopService: %w", err)
	}

	return &shopService{db: db}, nil
}

func (s *shopService) GetShop(ctx context.Context) (api.Response, error) {
	var shop database.PublicShopSelect
	err := s.db.DB.
		From("shop").
		Select("*").
		Single().
		Execute(&shop)

	if err != nil {
		log.Errorf("failed to fetch shop: %v", err)
		return nil, fmt.Errorf("failed to fetch shop information")
	}

	// Parse JSONB fields
	var branches []dto.Branch
	if shop.Branches != nil {
		branchesBytes, err := json.Marshal(shop.Branches)
		if err != nil {
			log.Errorf("failed to marshal branches: %v", err)
			return nil, fmt.Errorf("failed to parse shop branches")
		}
		if err := json.Unmarshal(branchesBytes, &branches); err != nil {
			log.Errorf("failed to unmarshal branches: %v", err)
			return nil, fmt.Errorf("failed to parse shop branches")
		}
	}

	var policies []dto.Policy
	if shop.Policies != nil {
		policiesBytes, err := json.Marshal(shop.Policies)
		if err != nil {
			log.Errorf("failed to marshal policies: %v", err)
			return nil, fmt.Errorf("failed to parse shop policies")
		}
		if err := json.Unmarshal(policiesBytes, &policies); err != nil {
			log.Errorf("failed to unmarshal policies: %v", err)
			return nil, fmt.Errorf("failed to parse shop policies")
		}
	}

	var contact []dto.Contact
	if shop.Contact != nil {
		contactBytes, err := json.Marshal(shop.Contact)
		if err != nil {
			log.Errorf("failed to marshal contact: %v", err)
			return nil, fmt.Errorf("failed to parse shop contact")
		}
		if err := json.Unmarshal(contactBytes, &contact); err != nil {
			log.Errorf("failed to unmarshal contact: %v", err)
			return nil, fmt.Errorf("failed to parse shop contact")
		}
	}

	response := dto.ShopResponse{
		Shop: dto.Shop{
			Id:        shop.Id,
			CreatedAt: parseTimeString(shop.CreatedAt),
			Name:      shop.Name,
			Hotline:   shop.Hotline,
			Zalo:      shop.Zalo,
			Branches:  branches,
			Policies:  policies,
			Contact:   contact,
		},
	}

	return response, nil
}

func (s *shopService) UpdateShop(ctx context.Context, req dto.ShopUpdate) (api.Response, error) {
	// Get the first shop ID
	var shop database.PublicShopSelect
	err := s.db.DB.
		From("shop").
		Select("id").
		Single().
		Execute(&shop)

	if err != nil {
		log.Errorf("failed to fetch shop for update: %v", err)
		return nil, fmt.Errorf("shop not found")
	}

	shopId := shop.Id

	// Prepare update data
	updateData := map[string]interface{}{
		"name":     req.Name,
		"hotline":  req.Hotline,
		"zalo":     req.Zalo,
		"branches": req.Branches,
		"policies": req.Policies,
		"contact":  req.Contact,
	}

	var updatedShops []database.PublicShopSelect
	err = s.db.DB.
		From("shop").
		Update(updateData).
		Eq("id", fmt.Sprintf("%d", shopId)).
		Execute(&updatedShops)

	if err != nil {
		log.Errorf("failed to update shop: %v", err)
		return nil, fmt.Errorf("failed to update shop information")
	}

	if len(updatedShops) == 0 {
		log.Error("no shop returned after update")
		return nil, fmt.Errorf("failed to update shop")
	}

	// Parse the updated shop data
	updatedShop := updatedShops[0]

	var branches []dto.Branch
	if updatedShop.Branches != nil {
		branchesBytes, _ := json.Marshal(updatedShop.Branches)
		json.Unmarshal(branchesBytes, &branches)
	}

	var policies []dto.Policy
	if updatedShop.Policies != nil {
		policiesBytes, _ := json.Marshal(updatedShop.Policies)
		json.Unmarshal(policiesBytes, &policies)
	}

	var contact []dto.Contact
	if updatedShop.Contact != nil {
		contactBytes, _ := json.Marshal(updatedShop.Contact)
		json.Unmarshal(contactBytes, &contact)
	}

	response := dto.ShopResponse{
		Shop: dto.Shop{
			Id:        updatedShop.Id,
			CreatedAt: parseTimeString(updatedShop.CreatedAt),
			Name:      updatedShop.Name,
			Hotline:   updatedShop.Hotline,
			Zalo:      updatedShop.Zalo,
			Branches:  branches,
			Policies:  policies,
			Contact:   contact,
		},
	}

	return response, nil
}

// parseTimeString parses a time string from the database
func parseTimeString(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
