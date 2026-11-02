package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	Auth AuthHandler
}

func New() *Handler {
	return &Handler{
		Auth: NewAuth(),
	}
}

type AuthHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	RefreshToken(c *gin.Context)
}
