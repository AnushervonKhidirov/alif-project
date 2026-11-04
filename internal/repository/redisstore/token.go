package redisstore

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenStore struct {
	rdb *redis.Client
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &tokenStore{
		rdb: client,
	}
}

func (s *tokenStore) Check(ctx *context.Context, userId int) error {
	key := s.generateTokenKey(userId)
	_, err := s.rdb.Exists(*ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *tokenStore) Get(ctx *context.Context, userId int) (string, error) {
	key := s.generateTokenKey(userId)
	token, err := s.rdb.Get(*ctx, key).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *tokenStore) Set(ctx *context.Context, userId int, token string, expiresAt time.Duration) error {
	key := s.generateTokenKey(userId)
	return s.rdb.Set(*ctx, key, token, expiresAt).Err()
}

func (s *tokenStore) Delete(ctx *context.Context, userId int) error {
	key := s.generateTokenKey(userId)
	s.rdb.Del(*ctx, key)
	return nil
}

func (*tokenStore) generateTokenKey(userId int) string {
	return "user_token_" + strconv.Itoa(userId)
}
