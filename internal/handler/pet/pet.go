package pet

import (
	"errors"
	"net/http"
	"pet/internal/dto/query"
	"pet/internal/service"
	"pet/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service service.Pet
}

func NewHandler(s service.Pet) *Handler {
	return &Handler{service: s}
}

func (h Handler) GetOne(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	pet, err := h.service.GetOne(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, pet)
}

func (h Handler) GetMany(c *gin.Context) {
	ctx := c.Request.Context()

	var filterQuery query.PetFilterMany

	err := c.ShouldBindQuery(&filterQuery)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	pets, err := h.service.GetMany(ctx, filterQuery)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, pets)
}

func (h Handler) Create(c *gin.Context) {}

func (h Handler) Update(c *gin.Context) {}

func (h Handler) Delete(c *gin.Context) {}

func (h Handler) MakeLost(c *gin.Context) {}

func (h Handler) UnmakeLost(c *gin.Context) {}

// my pet handlers
func (h Handler) UpdateMy(c *gin.Context) {}

func (h Handler) DeleteMy(c *gin.Context) {}

func (h Handler) MakeLostMy(c *gin.Context) {}

func (h Handler) UnmakeLostMy(c *gin.Context) {}
