package handler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/android-project-46group/core-api/config"
	pb "github.com/android-project-46group/protobuf/gen/go/protobuf"
	"google.golang.org/grpc"
)

type handler struct {
	pb.UnimplementedDownloadServer

	c config.Config
}

func New(c config.Config) pb.DownloadServer {
	h := &handler{
		UnimplementedDownloadServer: pb.UnimplementedDownloadServer{},
		c:                           c,
	}

	return h
}

func (h *handler) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
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
