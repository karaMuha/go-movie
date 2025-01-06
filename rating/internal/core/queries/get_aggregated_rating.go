package queries

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/ratingModel"
)

type GetAggregatedRatingQuery struct {
	ratingRepo driven.IRatingRepository
}

func NewGetAggregatedRatingQuery(ratingRepo driven.IRatingRepository) GetAggregatedRatingQuery {
	return GetAggregatedRatingQuery{
		ratingRepo: ratingRepo,
	}
}

func (q *GetAggregatedRatingQuery) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := q.ratingRepo.Load(ctx, recordID, recordType)
	if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	aggregatedRating := sum / float64(len(ratings))
	return aggregatedRating, nil
}
