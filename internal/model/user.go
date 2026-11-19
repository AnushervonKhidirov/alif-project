package model

import (
	"time"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
	FirstName    string
	LastName     *string
	Role         Role
	Pets         []Pet
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
