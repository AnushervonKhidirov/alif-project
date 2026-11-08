package handler

import (
	"pet/internal/handler/auth"
	pettype "pet/internal/handler/pet_type"
	"pet/internal/service"
)

type Handler struct {
	Auth    *auth.Handler
	PetType *pettype.Handler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth:    auth.NewHandler(s),
		PetType: pettype.NewHandler(s.PetType),
	}
}
