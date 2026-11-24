package query

import "pet/internal/model"

type Pagination struct {
	Limit  int `form:"take,default=20"`
	Offset int `form:"skip,default=0"`
}

type PetFilterMany struct {
	UserId      *int          `form:"userId"`
	Gender      *model.Gender `form:"gender"`
	MicrochipId *string       `form:"microchipId"`
	Lost        *bool         `form:"lost"`
	PetTypeId   *int          `form:"petTypeId"`
	BreedId     *int          `form:"breedId"`

	Pagination
}

type UserFilterMany struct {
	FirstName *string `form:"firstName"`
	LastName  *string `form:"lastName"`
	Email     *string `form:"email"`

	Pagination
}
