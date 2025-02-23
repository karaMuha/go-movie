package grpchandler

import (
	"context"
	"net/http"

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

	metadata, respErr := h.app.GetMetadata(ctx, req.MovieId)

	if respErr != nil {
		return &pb.GetMetadataResponse{
			ResponseStatus: metadataModel.RespErrToProto(respErr),
		}, nil
	}

	return &pb.GetMetadataResponse{
		Metadata: metadataModel.MetadataToProto(metadata),
		ResponseStatus: &pb.ResponseStatus{
			StatusCode: http.StatusOK,
		},
	}, nil
}

func (h *MetadataHandler) SubmitMetadata(ctx context.Context, req *pb.SubmitMetadataRequest) (*pb.SubmitMetadataResponse, error) {
	cmd := metadataModel.Metadata{
		Title:       req.Metadata.Title,
		Description: req.Metadata.Description,
		Director:    req.Metadata.Director,
	}

	metadata, respErr := h.app.CreateMetadata(ctx, &cmd)
	if respErr != nil {
		return &pb.SubmitMetadataResponse{
			ResponseStatus: metadataModel.RespErrToProto(respErr),
		}, nil
	}

	return &pb.SubmitMetadataResponse{
		Metadata: metadataModel.MetadataToProto(metadata),
		ResponseStatus: &pb.ResponseStatus{
			StatusCode: http.StatusCreated,
		},
	}, nil
}
