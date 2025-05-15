package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/log"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
	"net/http"
	"strings"
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
		return nil, fmt.Errorf("failed to initialize AuthService: %w", err)
	}
	return &authService{db: db}, nil
}

func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (api.Response, error) {
	if req.Password == "" || (req.Email == "" && req.Username == "") {
		return nil, fmt.Errorf("email/username and password are required")
	}

	// If login with username, you must first find the user's email
	email := req.Email
	if email == "" {
		var users []dto.User
		err := a.db.DB.From("users").Select("*").Eq("username", req.Username).Execute(&users)
		if err != nil || len(users) == 0 {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		log.Info(users[0].Username)
		email = users[0].Email
		log.Info(email)
	}

	session, err := a.db.Auth.SignIn(ctx, supabase.UserCredentials{
		Email:    email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	resp := dto.AuthResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}

	return api.Success(resp), nil
}

func (a *authService) Refresh(ctx context.Context, req dto.RefreshTokenRequest) (api.Response, error) {
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	authDetails, err := a.db.Auth.RefreshUser(ctx, "", req.RefreshToken)
	if err != nil {
		log.Errorf("failed to refresh session: %v", err)
		return nil, fmt.Errorf("failed to refresh token")
	}

	resp := dto.AuthResponse{
		AccessToken:  authDetails.AccessToken,
		RefreshToken: authDetails.RefreshToken,
	}

	return api.Success(resp), nil
}
func (a *authService) GetUserInfo(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := a.db.Auth.User(ctx, tokenString)
		if err != nil {
			return echo.ErrUnauthorized
		}
		return c.JSON(http.StatusOK, user)
	}
}
