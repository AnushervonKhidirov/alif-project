package pettype

import (
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.PetType
}

func NewHandler(s service.PetType) *Handler {
	return &Handler{service: s}
}

func (h Handler) GetOne(c *gin.Context) {
	// ctx := c.Request.Context()
}

func (h Handler) GetMany(c *gin.Context) {
	// ctx := c.Request.Context()
}

func (h Handler) Create(c *gin.Context) {
	// ctx := c.Request.Context()
}

func (h Handler) Update(c *gin.Context) {
	// ctx := c.Request.Context()
}

func (h Handler) Delete(c *gin.Context) {
	// ctx := c.Request.Context()
}
