package main

import (
	"context"

	"go-service/api"
	config "go-service/utils/configs"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	router := api.GetRouter(ctx)
	router.Run(cfg.Server.Address)
}
