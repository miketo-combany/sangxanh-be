package repository

import (
	"SangXanh/pkg/model"
	"github.com/samber/do/v2"
)

type UserRepository interface {
	Collection[*model.User]
}

type userRepository struct {
	collection[*model.User]
}

func NewUserRepository(di do.Injector) (UserRepository, error) {
	return &userRepository{
		collection: collection[*model.User]{
			db: newCollection(di, model.User{}.CollectionName()),
		},
	}, nil
}
