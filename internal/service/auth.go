package service

import "pet/internal/dto"

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (*authService) SignUp() (*dto.TokenResponse, error) {
	return nil, nil
}

func (*authService) SignIn() (*dto.TokenResponse, error) {
	return nil, nil
}

func (*authService) RefreshToken() (*dto.TokenResponse, error) {
	return nil, nil
}

func (*authService) SignOut() error {
	return nil
}
