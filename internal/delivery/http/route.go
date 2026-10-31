package route

import (
	"pet/internal/handler"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Engine *gin.Engine

	auth *AuthRoute
}

func New(handler *handler.Handler) *Route {
	engine := gin.Default()

	return &Route{
		Engine: engine,
		auth:   NewAuth(handler.Auth),
	}
}

func (r *Route) Init() {
	group := r.Engine.Group("api")

	{
		v1 := group.Group("v1")

		r.auth.Init(v1)
	}
}
