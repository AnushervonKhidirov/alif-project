package user

import (
	"pet/internal/handler/user"
	"pet/internal/middleware"
	"pet/internal/model"
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Route struct {
	h *user.Handler
	s service.TokenAccessVerifier
}

func NewRoute(handler *user.Handler, tokenService service.TokenAccessVerifier) *Route {
	return &Route{h: handler, s: tokenService}
}

func (r *Route) Init(g *gin.RouterGroup) {
	group := g.Group("/user")

	{
		{
			me := group.Group("/me")
			me.GET("/", middleware.Auth(r.s, nil), r.h.GetMe)
			me.PATCH("/", middleware.Auth(r.s, nil), r.h.UpdateMe)
			me.DELETE("/", middleware.Auth(r.s, nil), r.h.DeleteMe)
		}

		group.GET("/", middleware.Auth(r.s, []model.Role{model.RoleAdmin}), r.h.GetMany)
		group.PATCH("/change-password", middleware.Auth(r.s, nil), r.h.ChangePassword)
		group.GET("/:id", r.h.GetOne)
	}
}
