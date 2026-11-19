package config

import (
	"pet/pkg/database/postgres"
	"pet/pkg/database/redis"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppConf struct {
	Env        string          `env:"APP_ENV" env-default:"dev"`
	HttpServer HttpServer      `env-required:"true"`
	Database   postgres.Config `env-required:"true"`
	Redis      redis.Config    `env-required:"true"`
	Token      Token           `env-required:"true"`
}

type HttpServer struct {
	Port        string `env:"HTTP_PORT" env-default:"8080"`
	Host        string `env:"HTTP_HOST" env-default:"localhost"`
	Timeout     string `env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout string `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type Token struct {
	AccessSecretKey   string        `env:"TOKEN_ACCESS_KEY" env-required:"true"`
	RefreshSecretKey  string        `env:"TOKEN_REFRESH_KEY" env-required:"true"`
	AccessExpiration  time.Duration `env:"TOKEN_ACCESS_EXPIRATION" env-required:"true"`
	RefreshExpiration time.Duration `env:"TOKEN_REFRESH_EXPIRATION" env-required:"true"`
}

func InitConfig() (*AppConf, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var cfg AppConf
	err = cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
