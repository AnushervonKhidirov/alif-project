package model

import "time"

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type Pet struct {
	Id          int
	UserId      int
	Name        string
	About       *string
	Gender      *Gender
	Birthday    *time.Time
	MicrochipId *string
	Lost        bool
	PetTypeId   int
	PetType     PetType
	BreedId     *int
	Breed       *Breed
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
