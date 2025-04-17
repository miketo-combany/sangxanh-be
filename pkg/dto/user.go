package dto

import (
	"SangXanh/pkg/common/query"
	"SangXanh/pkg/enum"
)

type User struct {
	Id           string            `json:"id"`
	Username     string            `json:"username"`
	Password     string            `json:"password"`
	Role         enum.Role         `json:"role"`
	Address      []Address         `json:"address"`
	BasicAddress string            `json:"basic_address"`
	Avatar       string            `json:"avatar"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Metadata     map[string]string `json:"metadata"`
}

type Address struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	AddressJson      string `json:"address_json"`
	IsDefaultAddress bool   `json:"is_default_address"`
}

type UserRegisterRequest struct {
	Username     string            `json:"username"`
	Password     string            `json:"password"`
	Email        string            `json:"email"`
	Phone        string            `json:"phone"`
	Avatar       string            `json:"avatar"`
	BasicAddress string            `json:"basic_address"`
	Metadata     map[string]string `json:"metadata"`
}

type UserUpdateRequest struct {
	Id           string            `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	Avatar       string            `json:"avatar"`
	Phone        string            `json:"phone"`
	BasicAddress string            `json:"basic_address"`
	Metadata     map[string]string `json:"metadata"`
}

type UserUpdateAddressRequest struct {
	Id      string    `json:"id"`
	Address []Address `json:"address"`
}

type ListUser struct {
	query.Pagination
}
