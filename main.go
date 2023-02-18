package main

import (
	"context"
	"log"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/handler"
	"github.com/android-project-46group/core-api/util/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to initialize configuration: ", err)
	}

	logger, closer, err := logger.NewFileLogger(cfg.LogPath, "ubuntu", cfg.Service)
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}

	defer closer()

	h := handler.New(cfg, logger)
	if err := handler.ServeGRPC(cfg.GrpcPort, h); err != nil {
		logger.Errorf(context.Background(), "failed to serve: %v", err)
	}
}
