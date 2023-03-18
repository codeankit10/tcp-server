package main

import (
	"context"
	"fmt"
	"tcp-server/internal/client"
	"tcp-server/internal/pkg/config"
)

func main() {
	fmt.Println("starting client")
	//loading configuration

	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error loading in config", err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)

	address := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = client.Run(ctx, address)
	if err != nil {
		fmt.Println("client error", err)
	}

}
