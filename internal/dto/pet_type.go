package dto

type PetType struct {
	Value string `json:"value" binding:"required"`
}
