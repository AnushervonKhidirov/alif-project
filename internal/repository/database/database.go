package database

import (
	"context"
	"pet/internal/config"
	"pet/internal/model"
	"pet/internal/repository/database/user"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	cfg            *config.Database
	UserRepository UserRepo
}

func NewDatabase(cfg *config.Database) *Database {
	return &Database{
		cfg: cfg,
	}
}

func (db Database) Connect() (*pgx.Conn, error) {
	connStr := "postgres://" + db.cfg.User + ":" + db.cfg.Password + "@" + db.cfg.Host + ":" + db.cfg.Port + "/" + db.cfg.DBName
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	db.UserRepository = user.NewRepository(conn)

	return conn, nil
}

type UserRepo interface {
	GetOne(ctx context.Context, filters user.UserFilter) (*model.User, error)
	GetMany(ctx context.Context, filters user.UserFilter) ([]model.User, error)
	Create(ctx context.Context, userArgs user.CreateUserArgs) (*model.User, error)
	Update(ctx context.Context, userId int, userArgs user.UpdateUserArgs) (*model.User, error)
	Delete(ctx context.Context, userId int) error
}
