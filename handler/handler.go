package handler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/usecase"
	"github.com/android-project-46group/core-api/util/logger"
	pb "github.com/android-project-46group/protobuf/gen/go/protobuf"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type handler struct {
	pb.UnimplementedDownloadServer

	usecase usecase.Usecase
	config  config.Config
	logger  logger.Logger
}

func New(config config.Config, logger logger.Logger, usecase usecase.Usecase) pb.DownloadServer {
	handler := &handler{
		UnimplementedDownloadServer: pb.UnimplementedDownloadServer{},

		config:  config,
		logger:  logger,
		usecase: usecase,
	}

	return handler
}

func (h *handler) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.Health")
	defer span.Finish()

	h.logger.Info(ctx, "health check")

	return &pb.HealthReply{
		Message: "{\"health\": \"ok\"}",
	}, nil
}

func ServeGRPC(port string, server pb.DownloadServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDownloadServer(s, server)

	//nolint:wrapcheck
	return s.Serve(lis)
}
