package queries

import (
	"context"
	"errors"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/movie/movieModel"
	"github.com/karaMuha/go-movie/rating/ratingModel"
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

func (q *GetMovieDetailsQuery) GetMovieDetails(ctx context.Context, movieID string) (*movieModel.MovieDetails, error) {
	metadata, err := q.metadataGateway.GetMetadata(ctx, movieID)
	if err != nil {
		return nil, err
	}

	rating, err := q.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(movieID), ratingModel.RecordTypeMovie)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	movieDetails := movieModel.MovieDetails{
		Rating:   rating,
		Metadata: *metadata,
	}

	return &movieDetails, nil
}
