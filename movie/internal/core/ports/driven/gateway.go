package driven

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	ratingModel "github.com/karaMuha/go-movie/rating/ratingModel"
)

type IMetadataGateway interface {
	GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, error)
}

type IRatingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
	SubmitRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error
}
