package dto

import (
	"errors"
	"pet/internal/model"
	"pet/pkg/nullable"
	"pet/pkg/validator"
	"time"
)

type UserResponse struct {
	Id        int           `json:"id"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  *string       `json:"lastName,omitempty"`
	Pets      []PetResponse `json:"pets,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func ToUserResponse(user *model.User) UserResponse {
	userResponse := UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Pets != nil {
		pets := make([]PetResponse, 0, len(user.Pets))

		for _, pet := range user.Pets {
			pets = append(pets, ToPetResponse(&pet))
		}

		userResponse.Pets = pets
	}

	return userResponse
}

type CreateUser struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"firstName"`
	LastName  *string `json:"lastName"`
}

func (dto CreateUser) Validate() error {
	if isEmail := validator.IsEmail(dto.Email); !isEmail {
		return errors.New("'email' should be an email")
	}

	if isString := validator.IsString(dto.Password); !isString {
		return errors.New("'password' should be a string")
	}

	if minSize := validator.MinLength(dto.Password, 6); !minSize {
		return errors.New("'password' should be greater or equal to 6 symbols")
	}

	if isString := validator.IsString(dto.FirstName); !isString {
		return errors.New("'firstName' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.FirstName); !isBlank {
		return errors.New("'firstName' should not be blank")
	}

	if dto.LastName != nil {
		if isString := validator.IsString(*dto.LastName); !isString {
			return errors.New("'lastName' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.LastName); !isBlank {
			return errors.New("'lastName' should not be blank")
		}
	}

	return nil
}

type UpdateUser struct {
	FirstName nullable.String `json:"firstName"`
	LastName  nullable.String `json:"lastName"`
}

func (dto UpdateUser) Validate() error {
	if dto.FirstName.Value != nil {
		if isString := validator.IsString(*dto.FirstName.Value); !isString {
			return errors.New("'firstName' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.FirstName.Value); !isBlank {
			return errors.New("'firstName' should not be blank")
		}
	}

	if dto.LastName.Value != nil {
		if isString := validator.IsString(*dto.LastName.Value); !isString {
			return errors.New("'lastName' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.LastName.Value); !isBlank {
			return errors.New("'lastName' should not be blank")
		}
	}

	return nil
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (dto ChangePassword) Validate() error {
	if isString := validator.IsString(dto.OldPassword); !isString {
		return errors.New("'oldPassword' should be a string")
	}

	if minSize := validator.MinLength(dto.OldPassword, 6); !minSize {
		return errors.New("'oldPassword' should be greater or equal to 6 symbols")
	}

	if isString := validator.IsString(dto.NewPassword); !isString {
		return errors.New("'newPassword' should be a string")
	}

	if minSize := validator.MinLength(dto.NewPassword, 6); !minSize {
		return errors.New("'newPassword' should be greater or equal to 6 symbols")
	}

	return nil
}
