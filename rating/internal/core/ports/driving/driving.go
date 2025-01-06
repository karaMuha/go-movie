package driving

import (
	"context"

	model "github.com/karaMuha/go-movie/rating/ratingModel"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SaveRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type IQueries interface {
	GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error)
}
