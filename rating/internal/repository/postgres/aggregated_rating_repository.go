package postgres_repo

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type AggregatedRatingRepository struct {
	db *sql.DB
}

func NewAggregatedRatingRepository(db *sql.DB) AggregatedRatingRepository {
	return AggregatedRatingRepository{
		db: db,
	}
}

var _ driven.IAggregatedRatingRepository = (*AggregatedRatingRepository)(nil)

func (r *AggregatedRatingRepository) Load(ctx context.Context, recordID string, recordType string) (*ratingmodel.AggregatedRating, *dtos.RespErr) {
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
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return &aggregatedRating, nil
}

func (r *AggregatedRatingRepository) Save(ctx context.Context, metadata *ratingmodel.AggregatedRating) *dtos.RespErr {
	query := `
		INSERT INTO aggregated_ratings (id, record_type, rating, amount_ratings)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, metadata.ID, metadata.RecordType, metadata.Rating, metadata.AmountRatings)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return nil
}

func (r *AggregatedRatingRepository) Update(ctx context.Context, aggregatedRating *ratingmodel.AggregatedRating) *dtos.RespErr {
	query := `
		UPDATE aggregated_ratings
		SET rating = $1, amount_ratings = $2
		WHERE id = $3 AND record_type = $4;
	`
	_, err := r.db.ExecContext(ctx, query, aggregatedRating.Rating, aggregatedRating.AmountRatings, aggregatedRating.ID, aggregatedRating.RecordType)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return nil
}
