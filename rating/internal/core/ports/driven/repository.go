package driven

import (
	"context"

	metadatamodel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IRatingRepository interface {
	Load(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) ([]*ratingmodel.Rating, *dtos.RespErr)
	Save(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr
}

type IAggregatedRatingRepository interface {
	Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) *dtos.RespErr
	Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, *dtos.RespErr)
	Update(ctx context.Context, aggregatedRating *ratingmodel.AggregatedRating) *dtos.RespErr
}

type IMetadataEventRepository interface {
	Save(ctx context.Context, event metadatamodel.MetadataEvent) *dtos.RespErr
	Load(ctx context.Context) ([]metadatamodel.MetadataEvent, *dtos.RespErr)
	Delete(ctx context.Context, ID, recordType string) *dtos.RespErr
}
