package handler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/android-project-46group/core-api/config"
	"github.com/android-project-46group/core-api/util/logger"
	pb "github.com/android-project-46group/protobuf/gen/go/protobuf"
	"google.golang.org/grpc"
)

type handler struct {
	pb.UnimplementedDownloadServer

	config config.Config
	logger logger.Logger
}

func New(c config.Config, l logger.Logger) pb.DownloadServer {
	handler := &handler{
		UnimplementedDownloadServer: pb.UnimplementedDownloadServer{},

		config: c,
		logger: l,
	}

	return handler
}

func (h *handler) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
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

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to Serve: %w", err)
	}

	return nil
}
