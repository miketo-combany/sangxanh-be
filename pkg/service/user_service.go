package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/enum"
	"SangXanh/pkg/log"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"time"

	"github.com/nedpals/supabase-go"
	"github.com/samber/do/v2"
)

type UserService interface {
	ListUser(ctx context.Context, user dto.ListUser) (api.Response, error)
	Register(ctx context.Context, req dto.UserRegisterRequest) (api.Response, error)
	UpdateUser(ctx context.Context, req dto.UserUpdateRequest) (api.Response, error)
	UpdateUserAddress(ctx context.Context, req dto.UserUpdateAddressRequest) (api.Response, error)
	GetUserById(ctx context.Context, id string) (api.Response, error) // ← NEW
	ChangePassword(ctx context.Context, req dto.ChangePassword) (api.Response, error)
	SendMagicLink(ctx context.Context, req dto.ResetPasswordRequest) (api.Response, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (api.Response, error)
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

func (s *userService) GetUserById(ctx context.Context, id string) (api.Response, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	var users []dto.UserInfo
	err := s.db.DB.
		From("users").
		Select("*").
		Eq("id", id).
		IsNull("deleted_at").
		Execute(&users)
	if err != nil {
		log.Errorf("failed to query user: %v", err)
		return nil, fmt.Errorf("failed to retrieve user")
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return api.Success(users[0]), nil
}

// List all users (excluding soft-deleted ones)
// ---------------------------------------------------------------------
// List + paginate users
// ---------------------------------------------------------------------

func (s *userService) ListUser(ctx context.Context, filter dto.ListUser) (api.Response, error) {
	// 1. how many records satisfy the filter?
	name := filter.Name
	total, err := s.countUsers(ctx, filter, name)
	if err != nil {
		return nil, err
	}

	// 2. fetch the current page
	var users []dto.User
	query := s.db.DB.
		From("users").
		// pull just the columns you really need
		Select("id,username,email,role,avatar,phone,basic_address,metadata,created_at,updated_at").
		LimitWithOffset(int(filter.Limit), int((filter.Page-1)*filter.Limit)).
		IsNull("deleted_at") // keep soft-deleted rows out

	// -----------------------------------------------------------------
	// apply the SAME filters you used in countUsers
	// -----------------------------------------------------------------
	if name != "" {
		encoded := url.QueryEscape("%" + name + "%")
		query = query.Like("username", encoded)
	}
	if filter.Role != "" {
		query = query.Eq("role", string(filter.Role))
	}
	if filter.Status != "" {
		query = query.Eq("status", string(filter.Status)) // example boolean / enum
	}

	if err := query.Execute(&users); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// 3. add pagination meta & return
	filter.Pagination.Total = int64(total)
	return api.SuccessPagination(users, &filter.Pagination), nil
}

// ---------------------------------------------------------------------
// helper: count the total number of rows for pagination
// ---------------------------------------------------------------------

func (s *userService) countUsers(ctx context.Context, filter dto.ListUser, name string) (int, error) {
	q := s.db.DB.
		From("users").
		Select("id").
		IsNull("deleted_at")

	// replicate *exactly* the same filters used in ListUser
	if name != "" {
		encoded := url.QueryEscape("%" + name + "%")
		q = q.Like("username", encoded)
	}
	if filter.Role != "" {
		q = q.Eq("role", string(filter.Role))
	}
	if filter.Status != "" {
		q = q.Eq("status", string(filter.Status))
	}

	var tmp []struct{} // we only care about the header that carries the count
	if err := q.Execute(&tmp); err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return len(tmp), nil
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

	userData := dto.UserRegisterData{
		FullName:     req.FullName,
		Username:     req.Username,
		Role:         enum.User,
		BasicAddress: req.BasicAddress,
		Avatar:       req.Avatar,
		Phone:        req.Phone,
		Metadata:     req.Metadata,
	}

	user := supabase.UserCredentials{
		Email:    req.Email,
		Password: req.Password,
		Data:     userData,
	}

	_, err = s.db.Auth.SignUp(ctx, user)
	if err != nil {
		log.Errorf("failed to insert user: %v", err)
		return nil, fmt.Errorf(err.Error())
	}

	return api.Success(user), nil
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
		"full_name":     req.FullName,
		"email":         req.Email,
		"avatar":        req.Avatar,
		"phone":         req.Phone,
		"basic_address": req.BasicAddress,
		"metadata":      req.Metadata,
		"updated_at":    time.Now(),
	}

	log.Info("===============")
	var updated []dto.User
	err = s.db.DB.From("users").Update(updateData).Eq("id", req.Id).Execute(&updated)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, fmt.Errorf("update failed")
	}

	return api.Success(updated[0]), nil
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

	updateData := map[string]interface{}{
		"address":    req.Address,
		"updated_at": time.Now(),
	}
	var updated []dto.User
	err = s.db.DB.From("users").Update(updateData).Eq("id", req.Id).Execute(&updated)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, fmt.Errorf("update failed")
	}

	return api.Success(updated[0]), nil
}

func (s *userService) ChangePassword(ctx context.Context, req dto.ChangePassword) (api.Response, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}
	userToken, ok := ctx.Value("token").(string)
	if !ok || userToken == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User Token not found in context")
	}

	var users []dto.User
	err := s.db.DB.From("users").Select("*").Eq("id", userID).Execute(&users)
	if err != nil {
		log.Errorf("failed to find user: %v", err)
		return nil, fmt.Errorf("user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(req.OldPassword)) != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User old password is not correct")
	}

	updateData := map[string]interface{}{
		"password":   req.NewPassword,
		"updated_at": time.Now(),
	}
	user, err := s.db.Auth.UpdateUser(ctx, userToken, updateData)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return nil, err
	}
	updateUser := map[string]interface{}{
		"password":   string(password),
		"updated_at": time.Now(),
	}
	err = s.db.DB.From("users").Update(updateUser).Eq("id", userID).Execute(nil)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, err
	}
	return api.Success(user), nil
}

func (s *userService) SendMagicLink(ctx context.Context, request dto.ResetPasswordRequest) (api.Response, error) {
	err := s.db.Auth.SendMagicLink(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	return api.Success("email sent successfully"), nil
}

func (s *userService) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (api.Response, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}
	userToken, ok := ctx.Value("token").(string)
	if !ok || userToken == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User Token not found in context")
	}
	var users []dto.User
	err := s.db.DB.From("users").Select("*").Eq("id", userID).Execute(&users)
	if err != nil {
		log.Errorf("failed to find user: %v", err)
		return nil, fmt.Errorf("user not found")
	}
	updateData := map[string]interface{}{
		"password":   request.NewPassword,
		"updated_at": time.Now(),
	}
	user, err := s.db.Auth.UpdateUser(ctx, userToken, updateData)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, err
	}
	err = s.db.DB.From("users").Update(updateData).Eq("id", userID).Execute(nil)
	if err != nil {
		log.Errorf("failed to update user: %v", err)
		return nil, err
	}

	return api.Success(user), nil
}
