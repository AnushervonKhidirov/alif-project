package service

import (
	"context"
	"errors"
	"pet/internal/config"
	"pet/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	accessSecretKey   []byte
	refreshSecretKey  []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewTokenService(cfg *config.Token) TokenService {
	return &tokenService{
		accessSecretKey:   []byte(cfg.AccessSecretKey),
		refreshSecretKey:  []byte(cfg.RefreshSecretKey),
		accessExpiration:  cfg.AccessExpiration,
		refreshExpiration: cfg.RefreshExpiration,
	}
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
	// claims := NewTokenClaims(id, role, s.refreshExpiration)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// tokenString, err := token.SignedString(s.refreshSecretKey)
	// if err != nil {
	// 	return "", err
	// }

	// expiresAt := time.Now().Add(s.refreshExpiration)
	// err = s.repository.SaveRefreshToken(ctx, tokenString, expiresAt)
	// if err != nil {
	// 	return "", err
	// }

	// return tokenString, nil
	return "", nil
}

func (s tokenService) VerifyRefresh(ctx *context.Context, tokenString string) (*TokenClaims, error) {
	// token, err := s.repository.GetRefreshToken(ctx, tokenString)
	// if err != nil {
	// 	return nil, errors.New(InvalidTokenError)
	// }

	// defer s.repository.DeleteRefreshToken(ctx, token.RefreshToken)

	// claims, err := s.parseToken(token.RefreshToken, s.refreshKeyFunc)
	// if err != nil {
	// 	return nil, err
	// }

	// return claims, nil
	return nil, nil
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
