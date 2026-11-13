package dto

import (
	"pet/internal/model"
	"pet/pkg/nullable"
	"time"
)

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
