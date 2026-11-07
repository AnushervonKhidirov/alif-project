package token

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	rdb *redis.Client
}

func NewStore(client *redis.Client) *TokenStore {
	return &TokenStore{
		rdb: client,
	}
}

func (s *TokenStore) Check(ctx context.Context, token string) error {
	_, err := s.rdb.Exists(ctx, token).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *TokenStore) Get(ctx context.Context, token string) (int, error) {
	userIdStr, err := s.rdb.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *TokenStore) Set(ctx context.Context, token string, userId int, expiresAt time.Duration) error {
	return s.rdb.Set(ctx, token, userId, expiresAt).Err()
}

func (s *TokenStore) Delete(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, token).Err()
}
