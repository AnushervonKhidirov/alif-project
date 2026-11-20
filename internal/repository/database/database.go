package database

import (
	"context"
	"pet/internal/model"
	"pet/internal/repository/database/breed"
	"pet/internal/repository/database/pet"
	"pet/internal/repository/database/pettype"
	"pet/internal/repository/database/user"
	"pet/pkg/database/postgres"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn              *pgx.Conn
	UserRepository    UserRepo
	PetRepository     PetRepo
	PetTypeRepository PetTypeRepo
	BreedRepository   BreedRepo
}

func New(ctx context.Context, cfg *postgres.Config) (*Database, error) {
	conn, err := postgres.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	db := Database{}

	db.UserRepository = user.NewRepository(conn)
	db.PetRepository = pet.NewRepository(conn)
	db.PetTypeRepository = pettype.NewRepository(conn)
	db.BreedRepository = breed.NewRepository(conn)

	return &db, nil
}

type UserRepo interface {
	GetOne(ctx context.Context, filters user.FilterOne) (*model.User, error)
	GetOneWithPet(ctx context.Context, filters user.FilterOne) (*model.User, error)
	GetMany(ctx context.Context, filters user.FilterMany) ([]model.User, error)
	Create(ctx context.Context, createArgs user.CreateUserArgs) (*model.User, error)
	Update(ctx context.Context, id int, updateArgs user.UpdateUserArgs) (*model.User, error)
	Delete(ctx context.Context, id int) error
	GetPassword(ctx context.Context, id int) (string, error)
}

type PetRepo interface {
	GetOne(ctx context.Context, id int) (*model.Pet, error)
	GetOneWithUser(ctx context.Context, id int) (*model.Pet, error)
	GetMany(ctx context.Context, filters pet.FilterMany) ([]model.Pet, error)
	GetManyWithUser(ctx context.Context, filters pet.FilterMany) ([]model.Pet, error)
	Create(ctx context.Context, createArgs pet.CreateArgs) (*model.Pet, error)
	Update(ctx context.Context, filters pet.FilterOne, updateArgs pet.UpdateArgs) (*model.Pet, error)
	Delete(ctx context.Context, filters pet.FilterOne) error
}

type PetTypeRepo interface {
	GetOne(ctx context.Context, id int) (*model.PetType, error)
	GetMany(ctx context.Context) ([]model.PetType, error)
	Create(ctx context.Context, value string) (*model.PetType, error)
	Update(ctx context.Context, id int, value string) (*model.PetType, error)
	Delete(ctx context.Context, id int) error
}

type BreedRepo interface {
	GetOne(ctx context.Context, id int) (*model.Breed, error)
	GetMany(ctx context.Context) ([]model.Breed, error)
	Create(ctx context.Context, breedArgs breed.CreateBreedArgs) (*model.Breed, error)
	Update(ctx context.Context, breedId int, updateBreedArgs breed.UpdateBreedArgs) (*model.Breed, error)
	Delete(ctx context.Context, petTypeId int) error
}
