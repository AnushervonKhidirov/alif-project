package route

import (
	"pet/internal/delivery/http/auth"
	pettype "pet/internal/delivery/http/pet_type"
	"pet/internal/handler"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Engine *gin.Engine

	auth    *auth.Route
	petType *pettype.Route
}

func New(handler *handler.Handler) *Route {
	engine := gin.Default()

	return &Route{
		Engine:  engine,
		auth:    auth.NewRoute(handler.Auth),
		petType: pettype.NewRoute(handler.PetType),
	}
}

func (r *Route) Init() {
	group := r.Engine.Group("/api")

	{
		v1 := group.Group("/v1")

		// Initializing routes
		r.auth.Init(v1)
	}
}
