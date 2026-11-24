package dto

import (
	"errors"
	"pet/internal/model"
	"pet/pkg/validator"
)

type CreateBreed struct {
	Value     string `json:"value"`
	PetTypeId int    `json:"petTypeId"`
}

func (dto CreateBreed) Validate() error {
	if isString := validator.IsString(dto.Value); !isString {
		return errors.New("'value' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.Value); !isBlank {
		return errors.New("'value' should not be blank")
	}

	if isNum := validator.IsNumber(dto.PetTypeId); !isNum {
		return errors.New("'petTypeId' should be a number")
	}

	if isPositive := validator.MinNum(dto.PetTypeId, 1); !isPositive {
		return errors.New("'petTypeId' should be a positive number")
	}

	return nil
}

type UpdateBreed struct {
	Value     *string `json:"value"`
	PetTypeId *int    `json:"petTypeId"`
}

func (dto UpdateBreed) Validate() error {
	if dto.Value != nil {
		if isString := validator.IsString(*dto.Value); !isString {
			return errors.New("'value' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.Value); !isBlank {
			return errors.New("'value' should not be blank")
		}
	}

	if dto.PetTypeId != nil {
		if isNum := validator.IsNumber(*dto.PetTypeId); !isNum {
			return errors.New("'petTypeId' should be a number")
		}

		if isPositive := validator.MinNum(*dto.PetTypeId, 1); !isPositive {
			return errors.New("'petTypeId' should be a positive number")
		}
	}

	return nil
}

type BreedResponse struct {
	Id        int    `json:"id"`
	Value     string `json:"value"`
	PetTypeId int    `json:"petTypeId"`
}

func ToBreedResponse(breed *model.Breed) BreedResponse {
	return BreedResponse{
		Id:        breed.Id,
		Value:     breed.Value,
		PetTypeId: breed.PetTypeId,
	}
}
