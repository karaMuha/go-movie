package postgres_repo

import (
	"context"
	"database/sql"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type MetadataRepository struct {
	db *sql.DB
}

func NewMetadataRepository(db *sql.DB) MetadataRepository {
	return MetadataRepository{
		db: db,
	}
}

var _ driven.IMetadatarepository = (*MetadataRepository)(nil)

func (m *MetadataRepository) Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, error) {
	return nil, nil
}

func (m *MetadataRepository) Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error {
	return nil
}
