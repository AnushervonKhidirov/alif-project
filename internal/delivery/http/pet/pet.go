package pet

import (
	"pet/internal/handler/pet"
	"pet/internal/middleware"
	"pet/internal/model"
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Route struct {
	h *pet.Handler
	s service.TokenAccessVerifier
}

func NewRoute(handler *pet.Handler, tokenService service.TokenAccessVerifier) *Route {
	return &Route{h: handler, s: tokenService}
}

func (r *Route) Init(g *gin.RouterGroup) {
	group := g.Group("/pet")

	{
		{
			my := group.Group("/my")

			my.POST("/", middleware.Auth(r.s, nil), r.h.CreateMy)
			my.POST("/lost/make/:id", middleware.Auth(r.s, nil), r.h.MakeLostMy)
			my.POST("/lost/unmake/:id", middleware.Auth(r.s, nil), r.h.UnmakeLostMy)

			my.PATCH("/:id", middleware.Auth(r.s, nil), r.h.UpdateMy)
			my.DELETE("/:id", middleware.Auth(r.s, nil), r.h.DeleteMy)
		}

		group.GET("/", r.h.GetMany)
		group.GET("/:id", r.h.GetOne)

		// only for admin
		group.POST("/", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Create)
		group.POST("/lost/make/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.MakeLost)
		group.POST("/lost/unmake/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.UnmakeLost)

		group.PATCH("/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Update)
		group.DELETE("/:id", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.Delete)
	}
}
