package driven

import (
	"context"

	metadatamodel "github.com/karaMuha/go-movie/metadata/pkg"
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

type IMetadataEventRepository interface {
	Save(ctx context.Context, event metadatamodel.MetadataEvent) error
	Load(ctx context.Context) ([]metadatamodel.MetadataEvent, error)
	Delete(ctx context.Context, ID, recordType string) error
}
