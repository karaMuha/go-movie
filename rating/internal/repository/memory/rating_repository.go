package memory

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/pkg"
)

type RatingRepository struct {
	data map[model.RecordType]map[model.RecordID][]*model.Rating
}

var _ driven.IRatingRepository = (*RatingRepository)(nil)

func NewRatingsRepository() RatingRepository {
	return RatingRepository{
		data: make(map[model.RecordType]map[model.RecordID][]*model.Rating),
	}
}

func (r *RatingRepository) Load(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, *dtos.RespErr) {
	if _, ok := r.data[recordType]; !ok {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusNotFound,
			StatusMessage: "Not Found",
		}
	}

	return r.data[recordType][recordID], nil
}

func (r *RatingRepository) Save(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) *dtos.RespErr {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]*model.Rating{}
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], rating)
	return nil
}
