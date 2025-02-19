package driven

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IMetadataGateway interface {
	GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, error)
	SubmitMetadata(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, error)
}

type IRatingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, error)
	SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}
