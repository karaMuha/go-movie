package grpcgateway

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/pkg/grpcutil"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type RatingGateway struct {
	registry discovery.Registry
}

var _ driven.IRatingGateway = (*RatingGateway)(nil)

func NewRatingGateway(registry discovery.Registry) RatingGateway {
	return RatingGateway{registry: registry}
}

func (g *RatingGateway) GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, *dtos.RespErr) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	defer conn.Close()

	client := pb.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &pb.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(ratingmodel.RecordTypeMovie),
	})
	if err != nil {
		return 0, 0, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	if resp.ResponseStatus.StatusCode > 299 {
		return 0.0, 0, &dtos.RespErr{
			StatusCode:    int(resp.ResponseStatus.StatusCode),
			StatusMessage: resp.ResponseStatus.Message,
		}
	}

	return resp.RatingValue, int(resp.AmountRating), nil
}

func (g *RatingGateway) SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	defer conn.Close()

	client := pb.NewRatingServiceClient(conn)
	resp, err := client.SubmitRating(ctx, &pb.SubmitRatingRequest{
		UserId:      rating.UserID,
		RecordId:    rating.RecordID,
		RecordType:  rating.RecordType,
		RatingValue: int32(rating.Value), // #nosec G115
	})
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	if resp.ResponseStatus.StatusCode > 299 {
		return &dtos.RespErr{
			StatusCode:    int(resp.ResponseStatus.StatusCode),
			StatusMessage: resp.ResponseStatus.Message,
		}
	}

	return nil
}
