package pet

import (
	"errors"
	"net/http"
	"pet/internal/dto"
	"pet/internal/dto/query"
	"pet/internal/handler/helper"
	"pet/internal/service"
	"pet/pkg/nullable"
	"pet/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	pet, err := h.service.GetOne(ctx, id, true)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.ToPetResponse(pet))
}

func (h Handler) GetMany(c *gin.Context) {
	ctx := c.Request.Context()

	var filterQuery query.PetFilterMany

	err := c.ShouldBindQuery(&filterQuery)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	pets, err := h.service.GetMany(ctx, filterQuery, true)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	petsResponse := make([]dto.PetResponse, 0, len(pets))

	for _, pet := range pets {
		petsResponse = append(petsResponse, dto.ToPetResponse(&pet))
	}

	response.OK(c, http.StatusOK, petsResponse)
}

func (h Handler) Create(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("userId param should be a number"))
		return
	}

	h.createPet(c, userId)
}

func (h Handler) Update(c *gin.Context) {
	h.updatePet(c, nil)
}

func (h Handler) Delete(c *gin.Context) {
	h.deletePet(c, nil)
}

func (h Handler) MakeLost(c *gin.Context) {
	h.lostPet(c, nil, true)
}

func (h Handler) UnmakeLost(c *gin.Context) {
	h.lostPet(c, nil, false)
}

// my pet handlers
func (h Handler) CreateMy(c *gin.Context) {
	userId, err := helper.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, err)
		return
	}

	h.createPet(c, userId)
}

func (h Handler) UpdateMy(c *gin.Context) {
	userId, err := helper.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, err)
		return
	}

	h.updatePet(c, &userId)
}

func (h Handler) DeleteMy(c *gin.Context) {
	userId, err := helper.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, err)
		return
	}

	h.deletePet(c, &userId)
}

func (h Handler) MakeLostMy(c *gin.Context) {
	userId, err := helper.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, err)
		return
	}

	h.lostPet(c, &userId, true)
}

func (h Handler) UnmakeLostMy(c *gin.Context) {
	userId, err := helper.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, err)
		return
	}

	h.lostPet(c, &userId, false)
}

// common handlers
func (h Handler) createPet(c *gin.Context, userId int) {
	ctx := c.Request.Context()
	var data dto.CreatePet
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	pet, err := h.service.Create(ctx, userId, data)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.InvalidTextRepresentation:
				response.Fail(c, http.StatusBadRequest, errors.New(pgErr.Message))
				return
			case pgerrcode.ForeignKeyViolation:
				switch {
				case strings.Contains(pgErr.Detail, "breed_id"):
					response.Fail(c, http.StatusNotFound, errors.New("breed not found"))
				case strings.Contains(pgErr.Detail, "pet_type_id"):
					response.Fail(c, http.StatusNotFound, errors.New("pet type not found"))
				default:
					response.Fail(c, http.StatusBadRequest, nil)
				}

				return
			}
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusCreated, dto.ToPetResponse(pet))
}

func (h Handler) updatePet(c *gin.Context, userId *int) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	var data dto.UpdatePet
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	pet, err := h.service.Update(ctx, id, userId, data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet not found"))
			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.InvalidTextRepresentation:
				response.Fail(c, http.StatusBadRequest, errors.New(pgErr.Message))
				return
			case pgerrcode.ForeignKeyViolation:
				switch {
				case strings.Contains(pgErr.Detail, "breed_id"):
					response.Fail(c, http.StatusNotFound, errors.New("breed not found"))
				case strings.Contains(pgErr.Detail, "pet_type_id"):
					response.Fail(c, http.StatusNotFound, errors.New("pet type not found"))
				default:
					response.Fail(c, http.StatusBadRequest, nil)
				}

				return
			}
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.ToPetResponse(pet))
}

func (h Handler) deletePet(c *gin.Context, userId *int) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	if err := h.service.Delete(ctx, id, userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}

func (h Handler) lostPet(c *gin.Context, userId *int, lost bool) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errors.New("id param should be a number"))
		return
	}

	data := dto.UpdatePet{
		Lost: nullable.Bool{Value: &lost, IsSet: true},
	}

	if _, err = h.service.Update(ctx, id, userId, data); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("pet not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}
