package route

import (
	"pet/internal/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	h handler.AuthHandler
}

func NewAuth(handler handler.AuthHandler) *AuthRoute {
	return &AuthRoute{h: handler}
}

func (r *AuthRoute) Init(g *gin.RouterGroup) {
	group := g.Group("/auth")

	{
		group.POST("/sign-up", r.h.SignUp)
		group.POST("/sign-in", r.h.SignIn)
		group.POST("/sign-out", r.h.SignOut)
		group.POST("/refresh-token", r.h.RefreshToken)
	}
}
