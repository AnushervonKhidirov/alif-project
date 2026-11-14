package dto

import (
	"pet/internal/model"
	"pet/pkg/nullable"
	"time"
)

type PetResponse struct {
	Id          int            `json:"id"`
	UserId      int            `json:"userId"`
	User        *model.User    `json:"user,omitempty"`
	Name        string         `json:"name"`
	About       *string        `json:"about"`
	Gender      *model.Gender  `json:"gender"`
	Birthday    *time.Time     `json:"birthday"`
	MicrochipId *string        `json:"microchipId"`
	Lost        bool           `json:"lost"`
	PetTypeId   int            `json:"petTypeId"`
	PetType     *model.PetType `json:"petType,omitempty"`
	BreedId     *int           `json:"breedId"`
	Breed       *model.Breed   `json:"breed,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

func ToPetResponse(pet *model.Pet) PetResponse {
	return PetResponse{
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
}

type CreatePet struct {
	UserId      int           `json:"userId" binding:"required"`
	Name        string        `json:"name" binding:"required"`
	About       *string       `json:"about"`
	Gender      *model.Gender `json:"gender"`
	Birthday    *time.Time    `json:"birthday"`
	MicrochipId *string       `json:"microchipId"`
	PetTypeId   int           `json:"petTypeId" binding:"required"`
	BreedId     *int          `json:"breedId"`
}

// TODO: add json tags and validation for nullable
type UpdatePet struct {
	Name        nullable.String
	About       nullable.String
	Gender      nullable.Nullable[model.Gender]
	Birthday    nullable.Time
	MicrochipId nullable.String
	PetTypeId   nullable.Int
	BreedId     nullable.Int
}
