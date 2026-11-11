package handler

import (
	"pet/internal/handler/auth"
	"pet/internal/handler/breed"
	pettype "pet/internal/handler/pet_type"
	"pet/internal/handler/user"
	"pet/internal/service"
)

type Handler struct {
	Auth    *auth.Handler
	User    *user.Handler
	PetType *pettype.Handler
	Breed   *breed.Handler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth:    auth.NewHandler(s),
		User:    user.NewHandler(s.User),
		PetType: pettype.NewHandler(s.PetType),
		Breed:   breed.NewHandler(s.Breed),
	}
}
