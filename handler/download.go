package handler

import (
	pb "github.com/android-project-46group/protobuf/gen/go/protobuf"
	"github.com/opentracing/opentracing-go"
)

func (h *handler) DownloadMembersZip(
	//nolint:nosnakecase
	req *pb.DownloadMembersZipRequest, stream pb.Download_DownloadMembersZipServer,
) error {
	ctx := stream.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.DownloadMembersZip")
	defer span.Finish()

	h.logger.Infof(ctx, "handler.DownloadMembersZip: group ", req.Group.String())

	return nil
}
