package redisstore

import (
	"context"
	"pet/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client     *redis.Client
	TokenStore TokenStore
}

func New(cfg *config.Redis) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})

	return &RedisDB{
		Client:     rdb,
		TokenStore: NewTokenStore(rdb),
	}
}

type TokenStore interface {
	Check(ctx *context.Context, userId int) error
	Get(ctx *context.Context, userId int) (string, error)
	Set(ctx *context.Context, userId int, token string, expiresAt time.Duration) error
	Delete(ctx *context.Context, userId int) error
}
