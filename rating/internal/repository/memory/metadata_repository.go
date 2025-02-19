package memory

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/domain"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type MetadataRepository struct {
	data map[ratingmodel.RecordType]map[ratingmodel.RecordID]*ratingmodel.AggregatedRating
}

func NewMetadataRepository() MetadataRepository {
	return MetadataRepository{
		data: make(map[ratingmodel.RecordType]map[ratingmodel.RecordID]*ratingmodel.AggregatedRating),
	}
}

var _ driven.IMetadatarepository = (*MetadataRepository)(nil)

func (r *MetadataRepository) Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error {
	if _, ok := r.data[ratingmodel.RecordType(metadata.RecordType)]; !ok {
		r.data[ratingmodel.RecordType(metadata.RecordType)] = map[ratingmodel.RecordID]*ratingmodel.AggregatedRating{}
	}

	r.data[ratingmodel.RecordType(metadata.RecordType)][ratingmodel.RecordID(metadata.ID)] = metadata
	return nil
}

func (r *MetadataRepository) Load(ctx context.Context, recordID string, recrodType string) (*ratingmodel.AggregatedRating, error) {
	aggregatedRating, ok := r.data[ratingmodel.RecordType(recrodType)][ratingmodel.RecordID(recordID)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return aggregatedRating, nil
}
