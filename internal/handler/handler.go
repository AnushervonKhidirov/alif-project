package handler

import (
	"pet/internal/handler/auth"
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Auth AuthHandler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth: auth.NewHandler(s),
	}
}

type AuthHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	RefreshToken(c *gin.Context)
}
