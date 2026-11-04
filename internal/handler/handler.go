package handler

import (
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Auth AuthHandler
}

func New(s *service.Service) *Handler {
	return &Handler{
		Auth: NewAuth(s),
	}
}

type AuthHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	RefreshToken(c *gin.Context)
}
