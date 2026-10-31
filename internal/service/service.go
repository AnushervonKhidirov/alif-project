package service

import (
	"pet/internal/dto"
	"pet/internal/model"
)

type Service struct {
	Auth      AuthService
	User      UserService
	Pet       PetService
	PetType   PetTypeService
	Breed     BreedService
	VetClinic VetClinicService
}

func New() *Service {
	return &Service{}
}

// Auth
type AuthService interface {
	SignUp() (*dto.TokenResponse, error)
	SignIn() (*dto.TokenResponse, error)
	RefreshToken() (*dto.TokenResponse, error)
	SignOut() error
}

// User
type UserService interface {
}

// Pet
type PetService interface {
}

// PetType
type PetTypeService interface {
	GetById() (*model.PetType, error)
	GetAll() ([]model.PetType, error)
	Create() (*model.PetType, error)
	Update() (*model.PetType, error)
	Delete() error
}

// Breed
type BreedService interface {
	GetById() (*model.Breed, error)
	GetAll() ([]model.Breed, error)
	Create() (*model.Breed, error)
	Update() (*model.Breed, error)
	Delete() error
}

// VetClinic
type VetClinicService interface {
	GetById() (*model.VetClinic, error)
	GetAll() ([]model.VetClinic, error)
	Create() (*model.VetClinic, error)
	Update() (*model.VetClinic, error)
	Delete() error
}
