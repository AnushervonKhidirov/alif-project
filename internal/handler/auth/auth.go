package auth

import (
	"fmt"
	"net/http"
	"pet/internal/dto"
	"pet/internal/service"
	"pet/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService  service.UserService
	tokenService service.TokenService
}

func NewHandler(s *service.Service) *AuthHandler {
	return &AuthHandler{
		userService:  s.User,
		tokenService: s.Token,
	}
}

func (h AuthHandler) SignUp(c *gin.Context) {
	var data dto.CreateUser
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	u, err := h.userService.Create(ctx, &data)
	if err != nil {
		// TODO: handler several error types
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.tokenService.GeneratePair(ctx, u.Id, u.Role)
	if err != nil {
		// TODO: handler several error types
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	response.OK(c, http.StatusCreated, dto.TokenResponse(*tokenPair))
}

func (h AuthHandler) SignIn(c *gin.Context) {
	var data dto.SignInRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}

func (h AuthHandler) SignOut(c *gin.Context) {
	var data dto.RefreshTokenRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}

func (h AuthHandler) RefreshToken(c *gin.Context) {
	var data dto.RefreshTokenRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("%#v", data)
}
