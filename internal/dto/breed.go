package dto

import "pet/internal/model"

type CreateBreed struct {
	Value     string `json:"value" binding:"required"`
	PetTypeId int    `json:"petTypeId" binding:"required"`
}

type UpdateBreed struct {
	Value     *string `json:"value"`
	PetTypeId *int    `json:"petTypeId"`
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
