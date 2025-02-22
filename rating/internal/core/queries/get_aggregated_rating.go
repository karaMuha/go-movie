package queries

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/pkg"
)

type GetAggregatedRatingQuery struct {
	ratingRepo   driven.IRatingRepository
	metadataRepo driven.IAggregatedRatingRepository
}

func NewGetAggregatedRatingQuery(ratingRepo driven.IRatingRepository, metadataRepo driven.IAggregatedRatingRepository) GetAggregatedRatingQuery {
	return GetAggregatedRatingQuery{
		ratingRepo:   ratingRepo,
		metadataRepo: metadataRepo,
	}
}

func (q *GetAggregatedRatingQuery) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, int, error) {
	aggregatedRating, err := q.metadataRepo.Load(ctx, string(recordID), string(recordType))
	if err != nil {
		return 0, 0, err
	}
	return aggregatedRating.Rating, aggregatedRating.AmountRatings, nil
}
