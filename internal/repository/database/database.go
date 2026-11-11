package database

import (
	"context"
	"pet/internal/config"
	"pet/internal/model"
	"pet/internal/repository/database/breed"
	pettype "pet/internal/repository/database/pet_type"
	"pet/internal/repository/database/user"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	cfg               *config.Database
	UserRepository    UserRepo
	PetTypeRepository PetTypeRepo
	BreedRepository   BreedRepo
}

func New(cfg *config.Database) *Database {
	return &Database{cfg: cfg}
}

func (db *Database) Connect(ctx context.Context) (*pgx.Conn, error) {
	connStr := "postgres://" + db.cfg.User + ":" + db.cfg.Password + "@" + db.cfg.Host + ":" + db.cfg.Port + "/" + db.cfg.DBName
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	db.UserRepository = user.NewRepository(conn)
	db.PetTypeRepository = pettype.NewRepository(conn)
	db.BreedRepository = breed.NewRepository(conn)

	return conn, nil
}

type UserRepo interface {
	GetOne(ctx context.Context, filters user.Filters) (*model.User, error)
	GetMany(ctx context.Context, filters user.Filters) ([]model.User, error)
	Create(ctx context.Context, createArgs user.CreateUserArgs) (*model.User, error)
	Update(ctx context.Context, id int, updateArgs user.UpdateUserArgs) (*model.User, error)
	Delete(ctx context.Context, id int) error
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
