package service

import (
	"context"
	"pet/internal/config"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/cache"
	"pet/internal/repository/database"
	petDb "pet/internal/repository/database/pet"
	userDb "pet/internal/repository/database/user"
	"pet/internal/service/breed"
	"pet/internal/service/pet"
	petType "pet/internal/service/pet_type"
	"pet/internal/service/token"
	"pet/internal/service/user"
)

type Service struct {
	User    User
	Token   Token
	Pet     Pet
	PetType PetType
	Breed   Breed
}

func New(cfg *config.AppConf, r *database.Database, cache *cache.Cache) *Service {
	return &Service{
		User:    user.NewService(r.UserRepository),
		Token:   token.NewService(&cfg.Token, cache.Token),
		Pet:     pet.NewService(r.PetRepository),
		PetType: petType.NewService(r.PetTypeRepository),
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
	GetOne(ctx context.Context, filters userDb.FilterOne) (*model.User, error)
	GetMany(ctx context.Context, filters userDb.FilterMany) ([]model.User, error)
	Create(ctx context.Context, data dto.CreateUser) (*model.User, error)
	Update(ctx context.Context, id int, data dto.UpdateUser) (*model.User, error)
	ChangePassword(ctx context.Context, id int, data dto.ChangePassword) error
	Delete(ctx context.Context, id int) error
}

// Pet
type Pet interface {
	GetOne(ctx context.Context, id int) (*model.Pet, error)
	GetMany(ctx context.Context, filters petDb.FilterMany) ([]model.Pet, error)
	CreateMy(ctx context.Context, userId int, data dto.CreatePet) (*model.Pet, error)
	UpdateMy(ctx context.Context, filters petDb.FilterOne, data dto.UpdatePet) (*model.Pet, error)
	DeleteMy(ctx context.Context, filters petDb.FilterOne) error
	MakeLostMy(ctx context.Context, filters petDb.FilterOne) error
	UnmakeLostMy(ctx context.Context, filters petDb.FilterOne) error
}

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
