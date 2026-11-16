package handler

import (
	"pet/internal/handler/auth"
	"pet/internal/handler/breed"
	"pet/internal/handler/pet"
	"pet/internal/handler/pettype"
	"pet/internal/handler/user"
	"pet/internal/service"
)

type Handler struct {
	Auth    *auth.Handler
	User    *user.Handler
	Pet     *pet.Handler
	PetType *pettype.Handler
	Breed   *breed.Handler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth:    auth.NewHandler(s),
		User:    user.NewHandler(s.User),
		Pet:     pet.NewHandler(s.Pet),
		PetType: pettype.NewHandler(s.PetType),
		Breed:   breed.NewHandler(s.Breed),
	}
}
