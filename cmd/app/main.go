package main

import (
	"fmt"
	"pet/internal/config"
	route "pet/internal/delivery/http"
	"pet/internal/handler"
	"pet/internal/repository/redisstore"
	"pet/internal/service"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)

	rdb := redisstore.New(&cfg.Redis)

	s := service.New(cfg, rdb.Client)
	h := handler.New(s)

	r := route.New(h)
	r.Init()
	r.Engine.Run(":8080")
}
