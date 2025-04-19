package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/enum"
	"SangXanh/pkg/log"
	"context"
	"fmt"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUser(ctx context.Context, user dto.ListUser) (api.Response, error)
	Register(ctx context.Context, req dto.UserRegisterRequest) (api.Response, error)
	UpdateUser(ctx context.Context, req dto.UserUpdateRequest) (api.Response, error)
	UpdateUserAddress(ctx context.Context, req dto.UserUpdateAddressRequest) (api.Response, error)
}

type userService struct {
	db *supabase.Client
}

func NewUserService(di do.Injector) (UserService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UserService: %w", err)
	}
	return &userService{db: db}, nil
}

// List all users (excluding soft-deleted ones)
func (s *userService) ListUser(ctx context.Context, user dto.ListUser) (api.Response, error) {
	var users []dto.User
	err := s.db.DB.From("users").
		Select("*").
		IsNull("deleted_at").
		Execute(&users)
	if err != nil {
		log.Errorf("failed to list users: %v", err)
		return nil, fmt.Errorf("failed to list users")
	}
	return api.Success(users), nil
}

// Register a new user with validation and password hashing
func (s *userService) Register(ctx context.Context, req dto.UserRegisterRequest) (api.Response, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, fmt.Errorf("username, password and email are required")
	}

	// Check if user already exists
	var existing []dto.User
	err := s.db.DB.From("users").Select("*").Eq("username", req.Username).Execute(&existing)
	if err != nil {
		log.Errorf("failed to check existing user: %v", err)
		return nil, fmt.Errorf("failed to register user")
	}
	if len(existing) > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("failed to hash password: %v", err)
		return nil, fmt.Errorf("internal error")
	}

	userData := dto.UserRegisterData{
		Username:     req.Username,
		Role:         enum.User,
		BasicAddress: req.BasicAddress,
		Avatar:       req.Avatar,
		Phone:        req.Phone,
		Metadata:     req.Metadata,
	}

	user := supabase.UserCredentials{
		Email:    req.Email,
		Password: string(hashedPassword),
		Data:     userData,
	}

	var inserted []dto.User
	err = s.db.DB.From("users").Insert(user).Execute(&inserted)
	if err != nil {
		log.Errorf("failed to insert user: %v", err)
		return nil, fmt.Errorf("failed to register user")
	}

	return api.Success("User registered successfully"), nil
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UserUpdateRequest) (api.Response, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Check if user exists
	var users []dto.User
	err := s.db.DB.From("users").Select("*").Eq("id", req.Id).Execute(&users)
	if err != nil {
		log.Errorf("failed to find user: %v", err)
		return nil, fmt.Errorf("user not found")
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	updateData := map[string]interface{}{
		"username":      req.Username,
		"email":         req.Email,
		"avatar":        req.Avatar,
		"phone":         req.Phone,
		"basic_address": req.BasicAddress,
		"metadata":      req.Metadata,
		"updated_at":    time.Now(),
	}

	var updated []dto.User
	err = s.db.DB.From("users").Update(updateData).Eq("id", req.Id).Execute(&updated)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, fmt.Errorf("update failed")
	}

	return api.Success("User profile updated successfully"), nil
}

func (s *userService) UpdateUserAddress(ctx context.Context, req dto.UserUpdateAddressRequest) (api.Response, error) {
	if req.Id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Check if user exists
	var users []dto.User
	err := s.db.DB.From("users").Select("*").Eq("id", req.Id).Execute(&users)
	if err != nil {
		log.Errorf("failed to find user: %v", err)
		return nil, fmt.Errorf("user not found")
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	var updated []dto.User
	err = s.db.DB.From("users").Update(req.Address).Eq("id", req.Id).Execute(&updated)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, fmt.Errorf("update failed")
	}

	return api.Success("User profile updated successfully"), nil
}
