package service

import (
	"context"
	"pet/internal/config"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/cache"
	"pet/internal/repository/database"
	userDb "pet/internal/repository/database/user"
	"pet/internal/service/breed"
	pettype "pet/internal/service/pet_type"
	"pet/internal/service/token"
	"pet/internal/service/user"
)

type Service struct {
	User      User
	Token     Token
	Pet       Pet
	PetType   PetType
	Breed     Breed
	VetClinic VetClinic
}

func New(cfg *config.AppConf, r *database.Database, cache *cache.Cache) *Service {
	return &Service{
		User:    user.NewService(r.UserRepository),
		Token:   token.NewService(&cfg.Token, cache.Token),
		PetType: pettype.NewService(r.PetTypeRepository),
		Breed:   breed.NewService(r.BreedRepository),
	}
}

// Auth
type Token interface {
	TokenAccessVerifier
	TokenRefreshVerifier
	TokenGenerator

	DeleteRefresh(ctx context.Context, tokenString string) error
}

type TokenAccessVerifier interface {
	VerifyAccess(tokenString string) (*token.Claims, error)
}

type TokenRefreshVerifier interface {
	VerifyRefresh(ctx context.Context, tokenString string) (*token.Claims, error)
}

type TokenGenerator interface {
	GeneratePair(ctx context.Context, id int, role model.Role) (*model.TokenPair, error)
	GenerateAccess(id int, role model.Role) (string, error)
	GenerateRefresh(ctx context.Context, id int, role model.Role) (string, error)
}

// User
type User interface {
	GetOne(ctx context.Context, filter userDb.Filters) (*model.User, error)
	GetMany(ctx context.Context, filter userDb.Filters) ([]model.User, error)
	Create(ctx context.Context, data dto.CreateUser) (*model.User, error)
	Update(ctx context.Context, id int, data dto.UpdateUser) (*model.User, error)
	ChangePassword(ctx context.Context, id int, data dto.ChangePassword) error
	Delete(ctx context.Context, id int) error
}

// Pet
type Pet interface{}

// PetType
type PetType interface {
	GetOne(ctx context.Context, id int) (*model.PetType, error)
	GetMany(ctx context.Context) ([]model.PetType, error)
	Create(ctx context.Context, data dto.PetType) (*model.PetType, error)
	Update(ctx context.Context, id int, data dto.PetType) (*model.PetType, error)
	Delete(ctx context.Context, id int) error
}

// Breed
type Breed interface {
	GetOne(ctx context.Context, id int) (*model.Breed, error)
	GetMany(ctx context.Context) ([]model.Breed, error)
	Create(ctx context.Context, data dto.CreateBreed) (*model.Breed, error)
	Update(ctx context.Context, id int, data dto.UpdateBreed) (*model.Breed, error)
	Delete(ctx context.Context, id int) error
}

// VetClinic
type VetClinic interface {
	GetOne(ctx context.Context, id int) (*model.VetClinic, error)
	GetMany(ctx context.Context) ([]model.VetClinic, error)
	Create(ctx context.Context) (*model.VetClinic, error)
	Update(ctx context.Context, id int) (*model.VetClinic, error)
	Delete(ctx context.Context, id int) error
}
