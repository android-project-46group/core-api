package handler

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/android-project-46group/protobuf/gen/go/protobuf"
	"github.com/opentracing/opentracing-go"
)

const (
	// 一度の stream.Send で送信する容量。
	streamDataCapacity = 100 * 1024 * 1024 // 100 MB
)

func (h *handler) DownloadMembersZip(
	//nolint:nosnakecase
	req *pb.DownloadMembersZipRequest, stream pb.Download_DownloadMembersZipServer,
) error {
	ctx := stream.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.DownloadMembersZip")
	defer span.Finish()

	h.logger.Infof(ctx, "handler.DownloadMembersZip: group ", req.Group.String())

	reader, err := h.usecase.DownloadMembersZip(ctx)
	if err != nil {
		return fmt.Errorf("failed to DownloadMembersZip: %w", err)
	}

	buf := make([]byte, streamDataCapacity)

	for {
		size, err := reader.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}

		err = stream.Send(&pb.DownloadMembersZipReply{
			Data: buf[:size],
		})
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}
	}

	return nil
}
