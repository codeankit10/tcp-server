package main

import (
	"context"
	"fmt"
	"math/rand"
	"tcp-server/internal/pkg/cache"
	"tcp-server/internal/pkg/clock"
	"tcp-server/internal/pkg/config"
	"tcp-server/internal/server"
	"time"
)

func main() {
	fmt.Println("starting server")
	//loading configuration

	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error loading in config", err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	ctx = context.WithValue(ctx, "clock", clock.SystemClock{})

	cacheInst, err := cache.NewRedisCache(ctx, configInst.CacheHost, configInst.CachePort)
	if err != nil {
		fmt.Println("error init cache", err)
		return
	}

	ctx = context.WithValue(ctx, "cache", cacheInst)
	rand.Seed(time.Now().UnixNano())
	serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = server.Run(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error", err)
	}

}
