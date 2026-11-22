package dto

import "pet/internal/model"

type PetType struct {
	Value string `json:"value" binding:"required"`
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
