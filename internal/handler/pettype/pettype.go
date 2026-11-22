package pettype

import (
	"errors"
	"net/http"
	"pet/internal/dto"
	"pet/internal/service"
	"pet/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service service.PetType
}

func NewHandler(s service.PetType) *Handler {
	return &Handler{service: s}
}

func (h Handler) GetOne(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	petType, err := h.service.GetOne(ctx, id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, errors.New("pet type not found"))
		return
	}

	response.OK(c, http.StatusOK, dto.ToPetTypeResponse(petType))
}

func (h Handler) GetMany(c *gin.Context) {
	ctx := c.Request.Context()

	petTypes, err := h.service.GetMany(ctx)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	petTypesResponse := make([]dto.PetTypeResponse, 0, len(petTypes))
	for _, petType := range petTypes {
		petTypesResponse = append(petTypesResponse, dto.ToPetTypeResponse(&petType))
	}

	response.OK(c, http.StatusOK, petTypesResponse)
}

func (h Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.PetType
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	petType, err := h.service.Create(ctx, data)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				response.Fail(c, http.StatusConflict, errors.New("pet type already exists"))
				return
			}
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusCreated, dto.ToPetTypeResponse(petType))
}

func (h Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	var data dto.PetType
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	petType, err := h.service.Update(ctx, id, data)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				response.Fail(c, http.StatusConflict, errors.New("pet type already exists"))
				return
			}
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.ToPetTypeResponse(petType))
}

func (h Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet type not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}
