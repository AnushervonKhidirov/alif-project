package dto

type CreateBreed struct {
	Value     string `json:"value" binding:"required"`
	PetTypeId int    `json:"petTypeId" binding:"required"`
}

type UpdateBreed struct {
	Value     *string `json:"value"`
	PetTypeId *int    `json:"petTypeId"`
}
