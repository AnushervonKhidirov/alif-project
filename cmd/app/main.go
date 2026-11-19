package main

import (
	"context"
	"log"
	"pet/internal/config"
	route "pet/internal/delivery/http"
	"pet/internal/handler"
	"pet/internal/repository/cache"
	"pet/internal/repository/database"
	"pet/internal/service"
)

func main() {
	ctx := context.Background()
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	// cache redis
	cache := cache.New(&cfg.Redis)
	defer cache.Client.Close()

	// postgres
	r, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Conn.Close(ctx)

	s := service.New(cfg, r, cache)
	h := handler.New(s)

	router := route.New(h, s.Token)
	router.Init()
	router.Engine.Run(":" + cfg.HttpServer.Port)
}
