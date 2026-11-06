package main

import (
	"fmt"
	"pet/internal/config"
	route "pet/internal/delivery/http"
	"pet/internal/handler"
	"pet/internal/repository/cache"
	"pet/internal/service"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)

	cache := cache.New(&cfg.Redis)
	defer cache.Client.Close()

	s := service.New(cfg, cache)
	h := handler.New(s)

	r := route.New(h)
	r.Init()
	r.Engine.Run(":" + cfg.HttpServer.Port)
}
