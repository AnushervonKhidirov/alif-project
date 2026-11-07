package cache

import (
	"context"
	"pet/internal/config"
	"pet/internal/repository/cache/token"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	Token  Token
}

func New(cfg *config.Redis) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})

	return &Cache{
		Client: rdb,
		Token:  token.NewStore(rdb),
	}
}

type Token interface {
	Check(ctx context.Context, token string) error
	Get(ctx context.Context, token string) (int, error)
	Set(ctx context.Context, token string, userId int, expiresAt time.Duration) error
	Delete(ctx context.Context, token string) error
}
