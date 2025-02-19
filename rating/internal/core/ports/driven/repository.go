package driven

import (
	"context"

	model "github.com/karaMuha/go-movie/rating/pkg"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IRatingRepository interface {
	Load(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error)
	Save(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type IMetadatarepository interface {
	Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error
	Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, error)
}
