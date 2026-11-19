package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
}

func Connect(ctx context.Context, cfg *Config) (*pgx.Conn, error) {
	connStr := "postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DBName
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
