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
	c config.Config

	pb.UnimplementedDownloadServer
}

func New(c config.Config) handler {
	h := handler{
		c: c,
	}

	return h
}

func (h *handler) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	return &pb.HealthReply{
		Message: "{\"health\", \"ok\"}",
	}, nil
}

func (h *handler) Serve() error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", h.c.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDownloadServer(s, h)

	return s.Serve(lis)
}
