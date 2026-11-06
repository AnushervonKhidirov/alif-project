package token

import (
	"context"
	"errors"
	"pet/internal/config"
	"pet/internal/model"
	"pet/internal/repository/cache"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const InvalidTokenError string = "invalid token"

type Claims struct {
	UserId   int
	UserRole model.Role
	jwt.RegisteredClaims
}

func NewClaims(id int, role model.Role, Expiration time.Duration) Claims {
	now := time.Now()

	return Claims{
		UserId:   id,
		UserRole: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(Expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

type TokenService struct {
	redisClient       cache.Token
	accessSecretKey   []byte
	refreshSecretKey  []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewService(cfg *config.Token, redisClient cache.Token) *TokenService {
	return &TokenService{
		redisClient:       redisClient,
		accessSecretKey:   []byte(cfg.AccessSecretKey),
		refreshSecretKey:  []byte(cfg.RefreshSecretKey),
		accessExpiration:  cfg.AccessExpiration,
		refreshExpiration: cfg.RefreshExpiration,
	}
}

func (s TokenService) GeneratePair(ctx context.Context, id int, role model.Role) (*model.TokenPair, error) {
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

func (s TokenService) GenerateAccess(id int, role model.Role) (string, error) {
	claims := NewClaims(id, role, s.accessExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.accessSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s TokenService) VerifyAccess(tokenString string) (*Claims, error) {
	claims, err := s.parseToken(tokenString, s.accessKeyFunc)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s TokenService) GenerateRefresh(ctx context.Context, id int, role model.Role) (string, error) {
	claims := NewClaims(id, role, s.refreshExpiration)
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

func (s TokenService) VerifyRefresh(ctx context.Context, tokenString string) (*Claims, error) {
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
func (TokenService) parseToken(tokenString string, keyFunc jwt.Keyfunc) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(InvalidTokenError)
}

func (s TokenService) accessKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(InvalidTokenError)
	}
	return s.accessSecretKey, nil
}

func (s TokenService) refreshKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(InvalidTokenError)
	}
	return s.refreshSecretKey, nil
}
