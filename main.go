package main

import (
	"log"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/handler"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to initialize configuration: ", err)
	}

	h := handler.New(cfg)
	if err := handler.ServeGRPC(cfg.GrpcPort, h); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
