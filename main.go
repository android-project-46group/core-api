package main

import (
	"context"
	"log"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/handler"
	db "github.com/android-project-46group/core-api/repository/database"
	"github.com/android-project-46group/core-api/repository/remote"
	"github.com/android-project-46group/core-api/usecase"
	"github.com/android-project-46group/core-api/util"
	"github.com/android-project-46group/core-api/util/logger"
	"github.com/opentracing/opentracing-go"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to initialize configuration: ", err)
	}

	logger, fileCloser, err := logger.NewFileLogger(cfg.LogPath, "ubuntu", cfg.Service)
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}

	defer fileCloser()

	tracer, traceCloser, err := util.NewJaegerTracer(cfg.Service)
	if err != nil {
		logger.Errorf(context.Background(), "cannot initialize jaeger tracer: ", err)
	}

	defer traceCloser.Close()
	opentracing.SetGlobalTracer(tracer)

	db, err := db.New(context.Background(), cfg)
	if err != nil {
		logger.Errorf(context.Background(), "failed to db.New: ", err)
	}
	client := remote.New()

	usecase := usecase.New(db, client)

	h := handler.New(cfg, logger, usecase)
	if err := handler.ServeGRPC(cfg.GrpcPort, h); err != nil {
		logger.Errorf(context.Background(), "failed to serve: %v", err)
	}
}
