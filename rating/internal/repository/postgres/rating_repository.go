package postgres_repo

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type RatingRepository struct {
	db *sql.DB
}

var _ driven.IRatingRepository = (*RatingRepository)(nil)

func NewRatingRepository(db *sql.DB) RatingRepository {
	return RatingRepository{
		db: db,
	}
}

func (r *RatingRepository) Load(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) ([]*ratingmodel.Rating, *dtos.RespErr) {
	query := `
		SELECT *
		FROM ratings
		WHERE record_id = $1
		AND record_type = $2
	`

	rows, err := r.db.QueryContext(ctx, query, string(recordID), string(recordType))
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	var ratingList []*ratingmodel.Rating
	for rows.Next() {
		var rating ratingmodel.Rating
		err := rows.Scan(
			&rating.RecordID,
			&rating.RecordType,
			&rating.UserID,
			&rating.Value,
		)
		if err != nil {
			return nil, &dtos.RespErr{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: err.Error(),
			}
		}

		ratingList = append(ratingList, &rating)
	}

	return ratingList, nil
}

func (r *RatingRepository) Save(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr {
	query := `
		INSERT INTO ratings (record_id, record_type, user_id, value)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, rating.RecordID, rating.RecordType, rating.UserID, rating.Value)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}
	return nil
}
