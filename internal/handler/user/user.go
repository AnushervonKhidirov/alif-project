package user

import (
	"errors"
	"net/http"
	"pet/internal/dto"
	"pet/internal/repository/database/user"
	"pet/internal/service"
	"pet/pkg/helper/request"
	"pet/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service service.User
}

func NewHandler(s service.User) *Handler {
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

	user, err := h.service.GetOne(ctx, user.FilterOne{Id: &id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.ToUserResponse(user))
}

func (h Handler) GetMany(c *gin.Context) {
	ctx := c.Request.Context()

	firstName, firstNameOk := c.GetQuery("firstName")
	lastName, lastNameOk := c.GetQuery("lastName")

	pagination, err := request.GetPaginationQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
	}

	filters := user.FilterMany{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}

	if firstNameOk {
		filters.FirstName = &firstName
	}

	if lastNameOk {
		filters.LastName = &lastName
	}

	users, err := h.service.GetMany(ctx, filters)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	var usersResponse = make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		usersResponse = append(usersResponse, dto.ToUserResponse(&user))
	}

	response.OK(c, http.StatusOK, usersResponse)
}

func (h Handler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := request.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.GetOne(ctx, user.FilterOne{Id: &id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, dto.ToUserResponse(user))
}

func (h Handler) UpdateMe(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := request.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	var data dto.UpdateUser
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Update(ctx, id, data)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, user)
}

func (h Handler) DeleteMe(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := request.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}

func (h Handler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := request.GetUserIdFromRequest(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	var data dto.ChangePassword
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	if err = h.service.ChangePassword(ctx, id, data); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		response.Fail(c, http.StatusInternalServerError, nil)
		return
	}

	response.OK(c, http.StatusOK, nil)
}
