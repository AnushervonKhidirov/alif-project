package main

import (
	route "pet/internal/delivery/http"
	"pet/internal/handler"
)

func main() {

	h := handler.New()

	r := route.New(h)
	r.Init()
	r.Engine.Run(":8080")
}
