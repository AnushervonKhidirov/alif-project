package dto

import "pet/pkg/optional"

type CreateUser struct {
	Email     string          `json:"email" binding:"required,email"`
	Password  string          `json:"password" binding:"required,min=6"`
	FirstName string          `json:"firstName" binding:"required"`
	LastName  optional.String `json:"lastName"`
}

type UpdateUser struct {
	FirstName optional.String `json:"firstName"`
	LastName  optional.String `json:"lastName"`
}

type ChangePassword struct {
	Password string `json:"password" binding:"required,min=6"`
}
