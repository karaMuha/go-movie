package driven

import (
	"context"

	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IRatingRepository interface {
	Load(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) ([]*ratingmodel.Rating, error)
	Save(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type IAggregatedRatingRepository interface {
	Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error
	Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, error)
}
