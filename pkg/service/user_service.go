package service

import (
	"SangXanh/pkg/common/api"
	"SangXanh/pkg/common/errors"
	"SangXanh/pkg/dto"
	"SangXanh/pkg/log"
	"SangXanh/pkg/model"
	"SangXanh/pkg/repository"
	"context"
	"github.com/samber/do/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"time"
)

type UserService interface {
	ListUser(ctx context.Context, req dto.ListUser) (api.Response, error)
	CreateUser(ctx context.Context, req dto.CreateUser) (api.Response, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(di do.Injector) (UserService, error) {
	return &userService{
		userRepo: do.MustInvoke[repository.UserRepository](di),
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUser) (api.Response, error) {
	_, err := u.userRepo.FindOne(ctx, bson.M{
		"email": req.Email,
	})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorw("error when call userRepo.FindOne", "error", err, "email", req.Email)
		return nil, err
	}
	if err == nil {
		return nil, errors.BadRequest("user with email %v already exists", req.Email)
	}
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
	if err = u.userRepo.Create(ctx, user); err != nil {
		log.Errorw("error when call userRepo.Create", "error", err, "user", user)
		return nil, err
	}
	return api.Success(user), nil
}

func (u *userService) ListUser(ctx context.Context, req dto.ListUser) (api.Response, error) {
	q := req.Query()
	p := &req.Pagination
	users, err := u.userRepo.FindMany(ctx, q, p)
	if err != nil {
		log.Errorw("error when call userRepo.FindMany", "error", err, "req", req)
		return nil, err
	}
	return api.SuccessPagination(users, p), nil
}
