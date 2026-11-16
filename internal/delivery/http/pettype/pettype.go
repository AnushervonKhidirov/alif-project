package pettype

import (
	"pet/internal/handler/pettype"
	"pet/internal/middleware"
	"pet/internal/model"
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Route struct {
	h *pettype.Handler
	s service.TokenAccessVerifier
}

func NewRoute(handler *pettype.Handler, tokenService service.TokenAccessVerifier) *Route {
	return &Route{h: handler, s: tokenService}
}

func (r *Route) Init(g *gin.RouterGroup) {
	group := g.Group("/pet-type")

	{
		group.GET("/", r.h.GetMany)
		group.GET("/:id", r.h.GetOne)
		group.POST("/", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Create)
		group.PATCH("/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Update)
		group.DELETE("/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Delete)
	}
}
