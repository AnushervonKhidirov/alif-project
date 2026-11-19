package dto

import (
	"pet/internal/model"
	"pet/pkg/nullable"
	"time"
)

type UserResponse struct {
	Id        int           `json:"id"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  *string       `json:"lastName,omitempty"`
	Role      model.Role    `json:"role"`
	Pets      []PetResponse `json:"pets,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func ToUserResponse(user *model.User) UserResponse {
	userResponse := UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Pets != nil {
		pets := make([]PetResponse, 0, len(user.Pets))

		for _, pet := range user.Pets {
			pets = append(pets, ToPetResponse(&pet))
		}

		userResponse.Pets = pets
	}

	return userResponse
}

type CreateUser struct {
	Email     string          `json:"email" binding:"required,email"`
	Password  string          `json:"password" binding:"required,min=6"`
	FirstName string          `json:"firstName" binding:"required"`
	LastName  nullable.String `json:"lastName"`
}

type UpdateUser struct {
	FirstName nullable.String `json:"firstName"`
	LastName  nullable.String `json:"lastName"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"`
	NewPassword string `json:"NewPassword" binding:"required,min=6"`
}
