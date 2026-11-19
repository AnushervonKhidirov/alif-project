package dto

import (
	"pet/internal/model"
	"pet/pkg/nullable"
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
	Name        string        `json:"name" binding:"required"`
	About       *string       `json:"about"`
	Gender      *model.Gender `json:"gender"`
	Birthday    *time.Time    `json:"birthday"`
	MicrochipId *string       `json:"microchipId"`
	PetTypeId   int           `json:"petTypeId" binding:"required"`
	BreedId     *int          `json:"breedId"`
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
