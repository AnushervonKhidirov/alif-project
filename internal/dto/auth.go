package dto

import (
	"errors"
	"pet/pkg/validator"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto SignInRequest) Validate() error {
	if isEmail := validator.IsEmail(dto.Email); !isEmail {
		return errors.New("'email' should be an email")
	}

	if isString := validator.IsString(dto.Password); !isString {
		return errors.New("'password' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.Password); !isBlank {
		return errors.New("'password' should not be blank")
	}

	if minSize := validator.MinLength(dto.Password, 6); !minSize {
		return errors.New("'password' should be greater or equal to 6 symbols")
	}

	return nil
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (dto RefreshTokenRequest) Validate() error {
	if isString := validator.IsString(dto.RefreshToken); !isString {
		return errors.New("'refreshToken' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.RefreshToken); !isBlank {
		return errors.New("'refreshToken' should not be blank")
	}

	return nil
}
