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

func (s *TokenStore) Check(ctx context.Context, userId int) error {
	key := s.generateTokenKey(userId)
	_, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *TokenStore) Get(ctx context.Context, userId int) (string, error) {
	key := s.generateTokenKey(userId)
	token, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenStore) Set(ctx context.Context, userId int, token string, expiresAt time.Duration) error {
	key := s.generateTokenKey(userId)
	return s.rdb.Set(ctx, key, token, expiresAt).Err()
}

func (s *TokenStore) Delete(ctx context.Context, userId int) error {
	key := s.generateTokenKey(userId)
	s.rdb.Del(ctx, key)
	return nil
}

func (*TokenStore) generateTokenKey(userId int) string {
	return "user_token_" + strconv.Itoa(userId)
}
