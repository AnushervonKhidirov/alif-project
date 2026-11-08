package auth

import (
	"pet/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

type Route struct {
	h *auth.Handler
}

func NewRoute(handler *auth.Handler) *Route {
	return &Route{h: handler}
}

func (r *Route) Init(g *gin.RouterGroup) {
	group := g.Group("/auth")

	{
		group.POST("/sign-up", r.h.SignUp)
		group.POST("/sign-in", r.h.SignIn)
		group.POST("/sign-out", r.h.SignOut)
		group.POST("/refresh-token", r.h.RefreshToken)
	}
}
