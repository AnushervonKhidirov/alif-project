package route

import (
	"pet/internal/delivery/http/auth"
	"pet/internal/delivery/http/breed"
	"pet/internal/delivery/http/pettype"
	"pet/internal/delivery/http/user"
	"pet/internal/handler"
	"pet/internal/service"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Engine *gin.Engine

	auth    *auth.Route
	user    *user.Route
	petType *pettype.Route
	breed   *breed.Route
}

func New(handler *handler.Handler, tokenService service.TokenAccessVerifier) *Route {
	engine := gin.Default()

	return &Route{
		Engine:  engine,
		auth:    auth.NewRoute(handler.Auth),
		user:    user.NewRoute(handler.User, tokenService),
		petType: pettype.NewRoute(handler.PetType, tokenService),
		breed:   breed.NewRoute(handler.Breed, tokenService),
	}
}

func (r *Route) Init() {
	group := r.Engine.Group("/api")

	{
		v1 := group.Group("/v1")

		// Initializing routes
		r.auth.Init(v1)
		r.user.Init(v1)
		r.petType.Init(v1)
		r.breed.Init(v1)
	}
}
