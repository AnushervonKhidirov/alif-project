package auth

import (
	"errors"
	"net/http"
	"pet/internal/dto"
	"pet/internal/repository/database/user"
	"pet/internal/service"
	"pet/pkg/password"
	"pet/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	userService  service.User
	tokenService service.Token
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		userService:  s.User,
		tokenService: s.Token,
	}
}

func (h Handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.CreateUser
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	u, err := h.userService.Create(ctx, data)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				response.Fail(c, http.StatusConflict, errors.New("user already exists"))
				return
			}
		}

		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.tokenService.GeneratePair(ctx, u.Id, u.Role)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusCreated, dto.TokenResponse(*tokenPair))
}

func (h Handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.SignInRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	u, err := h.userService.GetOne(ctx, user.Filters{Email: &data.Email})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("wrong email or password"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	if correct := password.CheckHash(data.Password, u.PasswordHash); !correct {
		response.Fail(c, http.StatusBadRequest, errors.New("wrong email or password"))
		return
	}

	tokenPair, err := h.tokenService.GeneratePair(ctx, u.Id, u.Role)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.TokenResponse(*tokenPair))
}

func (h Handler) SignOut(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.RefreshTokenRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	_, err = h.tokenService.VerifyRefresh(ctx, data.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err)
		return
	}

	h.tokenService.DeleteRefresh(ctx, data.RefreshToken)
}

func (h Handler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.RefreshTokenRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	claims, err := h.tokenService.VerifyRefresh(ctx, data.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err)
		return
	}

	h.tokenService.DeleteRefresh(ctx, data.RefreshToken)

	tokenPair, err := h.tokenService.GeneratePair(ctx, claims.UserId, claims.UserRole)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.TokenResponse(*tokenPair))
}
