package dto

import (
	"SangXanh/pkg/common/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"id"`
	CreatedAt  int64              `json:"created_at"`
	UpdatedAt  int64              `json:"updated_at"`
	Email      string             `json:"email"`
	Name       string             `json:"name"`
	GivenName  string             `json:"given_name"`
	FamilyName string             `json:"family_name"`
	Avatar     string             `json:"avatar"`
	Metadata   any                `json:"metadata"`
}

type CreateUser struct {
	Email      string `json:"email" validate:"required,email"`
	Name       string `json:"name" validate:"required"`
	GivenName  string `json:"given_name" validate:"required"`
	FamilyName string `json:"family_name" validate:"required"`
	Avatar     string `json:"avatar" validate:"required,url"`
	Metadata   any    `json:"metadata"`
}

type ListUser struct {
	query.Pagination
	Email string `json:"email" query:"email"`
}

func (req *ListUser) Query() bson.M {
	return query.Query().Like("email", req.Email).BSON()
}
