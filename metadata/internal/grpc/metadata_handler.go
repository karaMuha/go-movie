package grpchandler

import (
	"context"
	"errors"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataHandler struct {
	pb.UnimplementedMetadataServiceServer
	app driving.IApplication
}

func NewMetadataHandler(app driving.IApplication) MetadataHandler {
	return MetadataHandler{
		app: app,
	}
}

func (h *MetadataHandler) GetMetadata(ctx context.Context, req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty movie id")
	}

	metadata, err := h.app.GetMetadata(ctx, req.MovieId)

	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.GetMetadataResponse{
		Metadata: metadataModel.MetadataToProto(metadata),
	}, nil
}
