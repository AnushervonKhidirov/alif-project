package dto

import (
	"errors"
	"pet/internal/model"
	"pet/pkg/validator"
)

type PetType struct {
	Value string `json:"value"`
}

func (dto PetType) Validate() error {
	if isString := validator.IsString(dto.Value); !isString {
		return errors.New("'value' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.Value); !isBlank {
		return errors.New("'value' should not be blank")
	}

	return nil
}

type PetTypeResponse struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

func ToPetTypeResponse(petType *model.PetType) PetTypeResponse {
	return PetTypeResponse{
		Id:    petType.Id,
		Value: petType.Value,
	}
}
