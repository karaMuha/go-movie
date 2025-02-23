package queries

import (
	"context"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/movie/movieModel"
	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type GetMovieDetailsQuery struct {
	metadataGateway driven.IMetadataGateway
	ratingGateway   driven.IRatingGateway
}

func NewGetMovieDetailsQuery(metadataGateway driven.IMetadataGateway, ratingGateway driven.IRatingGateway) GetMovieDetailsQuery {
	return GetMovieDetailsQuery{
		metadataGateway: metadataGateway,
		ratingGateway:   ratingGateway,
	}
}

func (q *GetMovieDetailsQuery) GetMovieDetails(ctx context.Context, movieID string) (*movieModel.MovieDetails, *dtos.RespErr) {
	metadata, respErr := q.metadataGateway.GetMetadata(ctx, movieID)
	if respErr != nil {
		return nil, respErr
	}

	rating, amountRatings, respErr := q.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(movieID), ratingmodel.RecordTypeMovie)
	if respErr != nil {
		return nil, respErr
	}

	movieDetails := movieModel.MovieDetails{
		Rating:        rating,
		AmountRatings: amountRatings,
		Metadata:      *metadata,
	}

	return &movieDetails, nil
}
