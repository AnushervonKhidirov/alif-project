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
	userService  service.UserService
	tokenService service.TokenService
}

func NewAuth(s *service.Service) AuthHandler {
	return &authHandler{
		userService:  s.User,
		tokenService: s.Token,
	}
}

func (h *authHandler) SignUp(c *gin.Context) {
	var data dto.CreateUser
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	u, err := h.userService.Create(&ctx, &data)
	if err != nil {
		// TODO: handler several error types
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.tokenService.GeneratePair(&ctx, u.Id, u.Role)
	if err != nil {
		// TODO: handler several error types
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	response.OK(c, http.StatusCreated, dto.TokenResponse(*tokenPair))
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
