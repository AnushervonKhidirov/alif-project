package main

import (
	"fmt"
	"pet/internal/config"
	route "pet/internal/delivery/http"
	"pet/internal/handler"
)

func main() {
	conf, err := config.InitConfig()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(conf)

	h := handler.New()

	r := route.New(h)
	r.Init()
	r.Engine.Run(":8080")
}
