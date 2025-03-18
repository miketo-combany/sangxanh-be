package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/model"
	"context"
	"github.com/samber/do/v2"
	"time"
)

type UserService interface {
	ListUser(ctx context.Context, req dto.ListUser) (api.Response, error)
	CreateUser(ctx context.Context, req dto.CreateUser) (api.Response, error)
}

type userService struct {
}

func NewUserService(di do.Injector) (UserService, error) {
	return &userService{}, nil
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUser) (api.Response, error) {
	user := &model.User{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:         req.Email,
		Name:          req.Name,
		GivenName:     req.GivenName,
		FamilyName:    req.FamilyName,
		Avatar:        req.Avatar,
		Metadata:      nil,
		Organizations: nil,
	}
	return api.Success(user), nil
}

func (u *userService) ListUser(ctx context.Context, req dto.ListUser) (api.Response, error) {
	p := &req.Pagination
	var users []*model.User
	return api.SuccessPagination(users, p), nil
}
