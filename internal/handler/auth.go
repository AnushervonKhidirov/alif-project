package handler

import (
	"fmt"
	"net/http"
	"pet/internal/dto"
	"pet/internal/service"
	"pet/pkg/response"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewAuth() AuthHandler {
	return &authHandler{}
}

func (h *authHandler) SignUp(c *gin.Context) {
	var data dto.SignUpRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
	}

	fmt.Printf("%#v", data)
}

func (h *authHandler) SignIn(c *gin.Context) {
	var data dto.SignInRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}

func (h *authHandler) SignOut(c *gin.Context) {
	var data dto.RefreshTokenRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	var data dto.RefreshTokenRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}
