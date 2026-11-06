package service

import (
	"context"
	"pet/internal/config"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/cache"
	"pet/internal/service/token"
	"pet/internal/service/user"
)

type Service struct {
	User      UserService
	Token     TokenService
	Pet       PetService
	PetType   PetTypeService
	Breed     BreedService
	VetClinic VetClinicService
}

func New(cfg *config.AppConf, cache *cache.Cache) *Service {
	return &Service{
		User:  user.NewService(),
		Token: token.NewService(&cfg.Token, cache.Token),
	}
}

// Auth
type TokenService interface {
	GeneratePair(ctx *context.Context, id int, role model.Role) (*model.TokenPair, error)
	GenerateAccess(id int, role model.Role) (string, error)
	VerifyAccess(tokenString string) (*token.Claims, error)
	GenerateRefresh(ctx *context.Context, id int, role model.Role) (string, error)
	VerifyRefresh(ctx *context.Context, tokenString string) (*token.Claims, error)
}

// User
type UserService interface {
	GetById(ctx *context.Context, id int) (*model.User, error)
	GetByEmail(ctx *context.Context, email string) (*model.User, error)
	GetMany(ctx *context.Context) ([]*model.User, error)
	Create(ctx *context.Context, data *dto.CreateUser) (*model.User, error)
	Update(ctx *context.Context, id int, data *dto.UpdateUser) error
	ChangePassword(ctx *context.Context, id int, data *dto.ChangePassword) error
	Delete(ctx *context.Context, id int) error
}

// Pet
type PetService interface {
}

// PetType
type PetTypeService interface {
	GetById() (*model.PetType, error)
	GetMany() ([]model.PetType, error)
	Create() (*model.PetType, error)
	Update() (*model.PetType, error)
	Delete() error
}

// Breed
type BreedService interface {
	GetById() (*model.Breed, error)
	GetMany() ([]model.Breed, error)
	Create() (*model.Breed, error)
	Update() (*model.Breed, error)
	Delete() error
}

// VetClinic
type VetClinicService interface {
	GetById() (*model.VetClinic, error)
	GetMany() ([]model.VetClinic, error)
	Create() (*model.VetClinic, error)
	Update() (*model.VetClinic, error)
	Delete() error
}
