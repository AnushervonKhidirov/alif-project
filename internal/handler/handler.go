package handler

import (
	"pet/internal/handler/auth"
	"pet/internal/service"
)

type Handler struct {
	Auth *auth.Handler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth: auth.NewHandler(s),
	}
}
