package redis

import "github.com/redis/go-redis/v9"

type Config struct {
	Host     string `env:"REDIS_HOST" env-default:"redis"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

func Connect(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})
}
