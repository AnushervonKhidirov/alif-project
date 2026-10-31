package route

import (
	"pet/internal/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	h *handler.Auth
}

func NewAuth(handler *handler.Auth) *AuthRoute {
	return &AuthRoute{h: handler}
}

func (r *AuthRoute) Init(g *gin.RouterGroup) {
	g.POST("sign-up", r.h.SignUp)
	g.POST("sign-in", r.h.SignIn)
	g.POST("sign-out", r.h.SignOut)
	g.POST("refresh-token", r.h.RefreshToken)
}
