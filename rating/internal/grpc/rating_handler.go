package grpchandler

import (
	"context"
	"errors"

	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/rating/internal/core/domain"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/rating/ratingModel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RatingHandler struct {
	pb.UnimplementedRatingServiceServer
	app driving.IApplication
}

func NewRatingHandler(app driving.IApplication) RatingHandler {
	return RatingHandler{
		app: app,
	}
}

func (h *RatingHandler) GetAggregatedRating(ctx context.Context, req *pb.GetAggregatedRatingRequest) (*pb.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record id or empty record type")
	}

	rating, err := h.app.GetAggregatedRating(ctx, ratingModel.RecordID(req.RecordId), ratingModel.RecordType(req.RecordType))

	if errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.GetAggregatedRatingResponse{
		RatingValue: rating,
	}, nil
}

func (h *RatingHandler) SubmitRating(ctx context.Context, req *pb.SubmitRatingRequest) (*pb.SubmitRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record id or empty user id")
	}

	rating := &ratingModel.Rating{
		RecordID:   req.RecordId,
		RecordType: string(ratingModel.RecordTypeMovie),
		UserID:     req.UserId,
		Value:      int(req.RatingValue),
	}
	err := h.app.SubmitRating(ctx, ratingModel.RecordID(req.RecordId), ratingModel.RecordTypeMovie, rating)

	return &pb.SubmitRatingResponse{}, err

}
