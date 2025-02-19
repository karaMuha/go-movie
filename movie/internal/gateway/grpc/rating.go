package grpcgateway

import (
	"context"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
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

func (g *RatingGateway) GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	client := pb.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &pb.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(ratingmodel.RecordTypeMovie),
	})
	if err != nil {
		return 0, 0, err
	}

	return resp.RatingValue, int(resp.AmountRating), nil
}

func (g *RatingGateway) SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewRatingServiceClient(conn)
	_, err = client.SubmitRating(ctx, &pb.SubmitRatingRequest{
		UserId:      rating.UserID,
		RecordId:    rating.RecordID,
		RecordType:  rating.RecordType,
		RatingValue: int32(rating.Value),
	})
	if err != nil {
		return err
	}

	return nil
}
