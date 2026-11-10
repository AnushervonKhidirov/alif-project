package breed

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
	service service.Breed
}

func NewHandler(s service.Breed) *Handler {
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

	breed, err := h.service.GetOne(ctx, id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, errors.New("breed not found"))
		return
	}

	response.OK(c, http.StatusOK, breed)
}

func (h Handler) GetMany(c *gin.Context) {
	ctx := c.Request.Context()

	breeds, err := h.service.GetMany(ctx)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, breeds)
}

func (h Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var data dto.CreateBreed
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	breed, err := h.service.Create(ctx, data)
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

	response.OK(c, http.StatusCreated, breed)
}

func (h Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	var data dto.UpdateBreed
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	breed, err := h.service.Update(ctx, id, data)
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

	response.OK(c, http.StatusOK, breed)
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
			response.Fail(c, http.StatusNotFound, errors.New("breed not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}
