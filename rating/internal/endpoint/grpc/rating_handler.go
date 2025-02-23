package grpchandler

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
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

	rating, amountRatings, respErr := h.app.GetAggregatedRating(ctx, ratingmodel.RecordID(req.RecordId), ratingmodel.RecordType(req.RecordType))

	if respErr != nil {
		return &pb.GetAggregatedRatingResponse{
			ResponseStatus: ratingmodel.RespErrToProto(respErr),
		}, nil
	}

	return &pb.GetAggregatedRatingResponse{
		RatingValue:  rating,
		AmountRating: int32(amountRatings),
		ResponseStatus: &pb.ResponseStatus{
			StatusCode: http.StatusOK,
		},
	}, nil
}

func (h *RatingHandler) SubmitRating(ctx context.Context, req *pb.SubmitRatingRequest) (*pb.SubmitRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty record id or empty user id")
	}

	rating := &ratingmodel.Rating{
		RecordID:   req.RecordId,
		RecordType: string(ratingmodel.RecordTypeMovie),
		UserID:     req.UserId,
		Value:      int(req.RatingValue),
	}
	respErr := h.app.SubmitRating(ctx, ratingmodel.RecordID(req.RecordId), ratingmodel.RecordTypeMovie, rating)
	if respErr != nil {
		return &pb.SubmitRatingResponse{
			ResponseStatus: ratingmodel.RespErrToProto(respErr),
		}, nil
	}

	return &pb.SubmitRatingResponse{
		ResponseStatus: &pb.ResponseStatus{
			StatusCode: http.StatusCreated,
		},
	}, nil

}
