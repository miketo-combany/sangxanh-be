package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/util"
	"context"
	"fmt"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (api.Response, error)
	Refresh(ctx context.Context, req dto.RefreshTokenRequest) (api.Response, error)
}

type authService struct {
	db *supabase.Client
}

func NewAuthService(di do.Injector) (AuthService, error) {
	db, err := do.Invoke[*supabase.Client](di)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize UserService: %w", err)
	}
	return &authService{db: db}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (api.Response, error) {
	var users []dto.User
	err := s.db.DB.From("users").
		Select("*").
		Eq("username", req.Username).
		Execute(&users)
	if err != nil || len(users) == 0 {
		return nil, fmt.Errorf("invalid username or password")
	}

	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	accessToken, err := util.GenerateAccessToken(user.Id, string(user.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := util.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	return api.Success(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}), nil
}

func (s *authService) Refresh(ctx context.Context, req dto.RefreshTokenRequest) (api.Response, error) {
	claims, err := util.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	newAccessToken, err := util.GenerateAccessToken(claims.UserID, claims.UserRole)
	if err != nil {
		return nil, err
	}

	return api.Success(dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
	}), nil
}
