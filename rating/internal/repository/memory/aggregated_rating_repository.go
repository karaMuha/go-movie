package memory

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/domain"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type AggregatedRatingRepository struct {
	data map[ratingmodel.RecordType]map[ratingmodel.RecordID]*ratingmodel.AggregatedRating
}

func NewAggregatedRatingRepository() AggregatedRatingRepository {
	return AggregatedRatingRepository{
		data: make(map[ratingmodel.RecordType]map[ratingmodel.RecordID]*ratingmodel.AggregatedRating),
	}
}

var _ driven.IAggregatedRatingRepository = (*AggregatedRatingRepository)(nil)

func (r *AggregatedRatingRepository) Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error {
	if _, ok := r.data[ratingmodel.RecordType(metadata.RecordType)]; !ok {
		r.data[ratingmodel.RecordType(metadata.RecordType)] = map[ratingmodel.RecordID]*ratingmodel.AggregatedRating{}
	}

	r.data[ratingmodel.RecordType(metadata.RecordType)][ratingmodel.RecordID(metadata.ID)] = metadata
	return nil
}

func (r *AggregatedRatingRepository) Load(ctx context.Context, recordID string, recrodType string) (*ratingmodel.AggregatedRating, error) {
	aggregatedRating, ok := r.data[ratingmodel.RecordType(recrodType)][ratingmodel.RecordID(recordID)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return aggregatedRating, nil
}

func (r *AggregatedRatingRepository) Update(ctx context.Context, aggregatedRating *ratingmodel.AggregatedRating) error {
	return nil
}
