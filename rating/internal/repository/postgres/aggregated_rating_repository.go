package postgres_repo

import (
	"context"
	"database/sql"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type AggregatedRatingRepository struct {
	db *sql.DB
}

func NewAggregatedMetadataRepository(db *sql.DB) AggregatedRatingRepository {
	return AggregatedRatingRepository{
		db: db,
	}
}

var _ driven.IAggregatedRatingRepository = (*AggregatedRatingRepository)(nil)

func (r *AggregatedRatingRepository) Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, error) {
	query := `
		SELECT * 
		FROM aggregated_ratings
		WHERE id = $1 AND record_type = $2
	`
	row := r.db.QueryRowContext(ctx, query, recordID, recordType)

	var aggregatedRating ratingmodel.AggregatedRating
	if err := row.Scan(
		&aggregatedRating.ID,
		&aggregatedRating.RecordType,
		&aggregatedRating.Rating,
		&aggregatedRating.AmountRatings,
	); err != nil {
		return nil, err
	}
	return &aggregatedRating, nil
}

func (r *AggregatedRatingRepository) Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) error {
	query := `
		INSERT INTO aggregated_ratings (id, record_type, rating, amount_ratings)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, metadata.ID, metadata.RecordType, metadata.Rating, metadata.AmountRatings)
	return err
}
