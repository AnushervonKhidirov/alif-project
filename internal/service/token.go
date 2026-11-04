package service

import (
	"context"
	"errors"
	"pet/internal/config"
	"pet/internal/model"
	"pet/internal/repository/redisstore"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const InvalidTokenError string = "invalid token"

type TokenClaims struct {
	UserId   int
	UserRole model.Role
	jwt.RegisteredClaims
}

func NewTokenClaims(id int, role model.Role, Expiration time.Duration) TokenClaims {
	now := time.Now()

	return TokenClaims{
		UserId:   id,
		UserRole: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(Expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

type tokenService struct {
	redisClient       redisstore.TokenStore
	accessSecretKey   []byte
	refreshSecretKey  []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewTokenService(cfg *config.Token, redisClient *redis.Client) TokenService {
	return &tokenService{
		accessSecretKey:   []byte(cfg.AccessSecretKey),
		refreshSecretKey:  []byte(cfg.RefreshSecretKey),
		accessExpiration:  cfg.AccessExpiration,
		refreshExpiration: cfg.RefreshExpiration,
	}
}

func (s tokenService) GeneratePair(ctx *context.Context, id int, role model.Role) (*model.TokenPair, error) {
	accessToken, err := s.GenerateAccess(id, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefresh(ctx, id, role)
	if err != nil {
		return nil, err
	}

	tokenPair := &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func (s tokenService) GenerateAccess(id int, role model.Role) (string, error) {
	claims := NewTokenClaims(id, role, s.accessExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.accessSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s tokenService) VerifyAccess(tokenString string) (*TokenClaims, error) {
	claims, err := s.parseToken(tokenString, s.accessKeyFunc)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s tokenService) GenerateRefresh(ctx *context.Context, id int, role model.Role) (string, error) {
	claims := NewTokenClaims(id, role, s.refreshExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.refreshSecretKey)
	if err != nil {
		return "", err
	}

	err = s.redisClient.Set(ctx, id, tokenString, s.refreshExpiration)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s tokenService) VerifyRefresh(ctx *context.Context, tokenString string) (*TokenClaims, error) {
	claims, err := s.parseToken(tokenString, s.refreshKeyFunc)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Check(ctx, claims.UserId)
	if err != nil {
		return nil, errors.New(InvalidTokenError)
	}

	return claims, nil
}

// helpers
func (tokenService) parseToken(tokenString string, keyFunc jwt.Keyfunc) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(InvalidTokenError)
}

func (s tokenService) accessKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(InvalidTokenError)
	}
	return s.accessSecretKey, nil
}

func (s tokenService) refreshKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(InvalidTokenError)
	}
	return s.refreshSecretKey, nil
}
