package dto

import (
	"errors"
	"fmt"
	"pet/internal/constant"
	"pet/internal/model"
	"pet/pkg/nullable"
	"pet/pkg/validator"
	"time"
)

type PetResponse struct {
	Id          int            `json:"id"`
	UserId      int            `json:"userId"`
	User        *UserResponse  `json:"user,omitempty"`
	Name        string         `json:"name"`
	About       *string        `json:"about,omitempty"`
	Gender      *model.Gender  `json:"gender,omitempty"`
	Birthday    *time.Time     `json:"birthday,omitempty"`
	MicrochipId *string        `json:"microchipId,omitempty"`
	Lost        bool           `json:"lost"`
	PetTypeId   int            `json:"petTypeId"`
	PetType     *model.PetType `json:"petType,omitempty"`
	BreedId     *int           `json:"breedId,omitempty"`
	Breed       *model.Breed   `json:"breed,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func ToPetResponse(pet *model.Pet) PetResponse {
	petResponse := PetResponse{
		Id:          pet.Id,
		UserId:      pet.UserId,
		Name:        pet.Name,
		About:       pet.About,
		Gender:      pet.Gender,
		Birthday:    pet.Birthday,
		MicrochipId: pet.MicrochipId,
		Lost:        pet.Lost,
		PetTypeId:   pet.PetTypeId,
		PetType:     pet.PetType,
		BreedId:     pet.BreedId,
		Breed:       pet.Breed,
		CreatedAt:   pet.CreatedAt,
		UpdatedAt:   pet.UpdatedAt,
	}

	if pet.User != nil {
		user := ToUserResponse(pet.User)
		petResponse.User = &user
	}

	return petResponse
}

type CreatePet struct {
	Name        string        `json:"name"`
	About       *string       `json:"about"`
	Gender      *model.Gender `json:"gender"`
	Birthday    *time.Time    `json:"birthday"`
	MicrochipId *string       `json:"microchipId"`
	PetTypeId   int           `json:"petTypeId"`
	BreedId     *int          `json:"breedId"`
}

func (dto CreatePet) Validate() error {
	if isString := validator.IsString(dto.Name); !isString {
		return errors.New("'name' should be a string")
	}

	if isBlank := validator.IsBlankString(dto.Name); !isBlank {
		return errors.New("'name' should not be blank")
	}

	if dto.About != nil {
		if isString := validator.IsString(*dto.About); !isString {
			return errors.New("'about' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.About); !isBlank {
			return errors.New("'about' should not be blank")
		}
	}

	if dto.Gender != nil {
		if isString := validator.IsString(*dto.Gender); !isString {
			return fmt.Errorf("'gender' should be '%s' of '%s'", model.GenderMale, model.GenderFemale)
		}

		if *dto.Gender != model.GenderMale && *dto.Gender != model.GenderFemale {
			return fmt.Errorf("'gender' should be '%s' of '%s'", model.GenderMale, model.GenderFemale)
		}
	}

	if dto.Birthday != nil {
		if isValid := validator.IsValidDate(*dto.Birthday, constant.DateLayout); !isValid {
			return fmt.Errorf("'birthday' should be a valid date (example: %s)", constant.DateLayout)
		}
	}

	if dto.MicrochipId != nil {
		if isString := validator.IsString(*dto.MicrochipId); !isString {
			return errors.New("'microchipId' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.MicrochipId); !isBlank {
			return errors.New("'microchipId' should not be blank")
		}
	}

	if isNum := validator.IsNumber(dto.PetTypeId); !isNum {
		return errors.New("'petTypeId' should be a number")
	}

	if isPositive := validator.MinNum(dto.PetTypeId, 1); !isPositive {
		return errors.New("'petTypeId' should be a positive number")
	}

	if dto.BreedId != nil {
		if isNum := validator.IsNumber(*dto.BreedId); !isNum {
			return errors.New("'breedId' should be a number")
		}

		if isPositive := validator.MinNum(*dto.BreedId, 1); !isPositive {
			return errors.New("'breedId' should be a positive number")
		}
	}

	return nil
}

type UpdatePet struct {
	Name        nullable.String                 `json:"name"`
	About       nullable.String                 `json:"about"`
	Gender      nullable.Nullable[model.Gender] `json:"gender"`
	Birthday    nullable.Time                   `json:"birthday"`
	MicrochipId nullable.String                 `json:"microchipId"`
	PetTypeId   nullable.Int                    `json:"petTypeId"`
	BreedId     nullable.Int                    `json:"breedId"`
	Lost        nullable.Bool                   `json:"lost"`
}

func (dto UpdatePet) Validate() error {
	if dto.Name.Value != nil {
		if isString := validator.IsString(*dto.Name.Value); !isString {
			return errors.New("'name' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.Name.Value); !isBlank {
			return errors.New("'name' should not be blank")
		}
	}

	if dto.About.Value != nil {
		if isString := validator.IsString(*dto.About.Value); !isString {
			return errors.New("'about' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.About.Value); !isBlank {
			return errors.New("'about' should not be blank")
		}
	}

	if dto.Gender.Value != nil {
		if isString := validator.IsString(*dto.Gender.Value); !isString {
			return fmt.Errorf("'gender' should be '%s' of '%s'", model.GenderMale, model.GenderFemale)
		}

		if *dto.Gender.Value != model.GenderMale && *dto.Gender.Value != model.GenderFemale {
			return fmt.Errorf("'gender' should be '%s' of '%s'", model.GenderMale, model.GenderFemale)
		}
	}

	if dto.Birthday.Value != nil {
		if isValid := validator.IsValidDate(*dto.Birthday.Value, constant.DateLayout); !isValid {
			return fmt.Errorf("'birthday' should be a valid date (example: %s)", constant.DateLayout)
		}
	}

	if dto.MicrochipId.Value != nil {
		if isString := validator.IsString(*dto.MicrochipId.Value); !isString {
			return errors.New("'microchipId' should be a string")
		}

		if isBlank := validator.IsBlankString(*dto.MicrochipId.Value); !isBlank {
			return errors.New("'microchipId' should not be blank")
		}
	}

	if dto.PetTypeId.Value != nil {
		if isNum := validator.IsNumber(*dto.PetTypeId.Value); !isNum {
			return errors.New("'petTypeId' should be a number")
		}

		if isPositive := validator.MinNum(*dto.PetTypeId.Value, 1); !isPositive {
			return errors.New("'petTypeId' should be a positive number")
		}
	}

	if dto.BreedId.Value != nil {
		if isNum := validator.IsNumber(*dto.BreedId.Value); !isNum {
			return errors.New("'breedId' should be a number")
		}

		if isPositive := validator.MinNum(*dto.BreedId.Value, 1); !isPositive {
			return errors.New("'breedId' should be a positive number")
		}
	}

	return nil
}
