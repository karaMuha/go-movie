package driving

import (
	"context"

	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr
	SubmitMetadata(cmd *ratingmodel.AggregatedRating) *dtos.RespErr
}

type IQueries interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, *dtos.RespErr)
}
