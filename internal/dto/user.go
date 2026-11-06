package dto

import "pet/pkg/optional"

type CreateUser struct {
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	FirstName string          `json:"firstName"`
	LastName  optional.String `json:"lastName"`
}

type UpdateUser struct {
	FirstName optional.String `json:"firstName"`
	LastName  optional.String `json:"lastName"`
}

type ChangePassword struct {
	Password string `json:"password"`
}
