package service

import (
	"context"
	"pet/internal/config"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/cache"
	"pet/internal/repository/database"
	userDb "pet/internal/repository/database/user"
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

func New(cfg *config.AppConf, repo *database.Database, cache *cache.Cache) *Service {
	return &Service{
		User:  user.NewService(repo.UserRepository),
		Token: token.NewService(&cfg.Token, cache.Token),
	}
}

// Auth
type TokenService interface {
	GeneratePair(ctx context.Context, id int, role model.Role) (*model.TokenPair, error)
	GenerateAccess(id int, role model.Role) (string, error)
	VerifyAccess(tokenString string) (*token.Claims, error)
	GenerateRefresh(ctx context.Context, id int, role model.Role) (string, error)
	VerifyRefresh(ctx context.Context, tokenString string) (*token.Claims, error)
	DeleteRefresh(ctx context.Context, tokenString string) error
}

// User
type UserService interface {
	GetOne(ctx context.Context, filter userDb.Filters) (*model.User, error)
	GetMany(ctx context.Context, filter userDb.Filters) ([]model.User, error)
	Create(ctx context.Context, data *dto.CreateUser) (*model.User, error)
	Update(ctx context.Context, id int, data *dto.UpdateUser) (*model.User, error)
	ChangePassword(ctx context.Context, id int, newPasswordHash string) error
	Delete(ctx context.Context, id int) error
}

// Pet
type PetService interface {
}

// PetType
type PetTypeService interface {
	GetOne() (*model.PetType, error)
	GetMany() ([]model.PetType, error)
	Create() (*model.PetType, error)
	Update() (*model.PetType, error)
	Delete() error
}

// Breed
type BreedService interface {
	GetOne() (*model.Breed, error)
	GetMany() ([]model.Breed, error)
	Create() (*model.Breed, error)
	Update() (*model.Breed, error)
	Delete() error
}

// VetClinic
type VetClinicService interface {
	GetOne() (*model.VetClinic, error)
	GetMany() ([]model.VetClinic, error)
	Create() (*model.VetClinic, error)
	Update() (*model.VetClinic, error)
	Delete() error
}
