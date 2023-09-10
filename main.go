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

	defer func() {
		if err := fileCloser(); err != nil {
			logger.Warnf(context.Background(), "failed to fileCloser: ", err)
		}
	}()

	tracer, traceCloser, err := util.NewJaegerTracer(cfg.Service)
	if err != nil {
		logger.Errorf(context.Background(), "cannot initialize jaeger tracer: ", err)
	}

	defer func() {
		if err := traceCloser.Close(); err != nil {
			logger.Warnf(context.Background(), "failed to traceCloser: ", err)
		}
	}()
	opentracing.SetGlobalTracer(tracer)

	database, err := db.New(context.Background(), cfg, logger)
	if err != nil {
		logger.Errorf(context.Background(), "failed to db.New: ", err)
	}

	client := remote.New()
	usecase := usecase.New(database, client, logger)

	h := handler.New(cfg, logger, usecase)
	if err := handler.ServeGRPC(cfg.GrpcPort, h); err != nil {
		logger.Errorf(context.Background(), "failed to serve: %v", err)
	}
}
