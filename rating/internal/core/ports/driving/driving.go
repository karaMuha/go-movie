package driving

import (
	"context"

	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
	SubmitMetadata(cmd *ratingmodel.AggregatedRating) error
}

type IQueries interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, error)
}
