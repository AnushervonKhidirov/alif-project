package handler

import "github.com/gin-gonic/gin"

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (h *Auth) SignUp(ctx *gin.Context) {}

func (h *Auth) SignIn(ctx *gin.Context) {}

func (h *Auth) SignOut(ctx *gin.Context) {}

func (h *Auth) RefreshToken(ctx *gin.Context) {}
