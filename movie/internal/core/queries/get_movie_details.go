package queries

import (
	"context"
	"strings"

	"github.com/karaMuha/go-movie/movie/internal/core/domain"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/movie/movieModel"
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

func (q *GetMovieDetailsQuery) GetMovieDetails(ctx context.Context, movieID string) (*movieModel.MovieDetails, error) {
	metadata, err := q.metadataGateway.GetMetadata(ctx, movieID)
	if strings.Contains(err.Error(), "NotFound") {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rating, err := q.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(movieID), ratingmodel.RecordTypeMovie)
	if strings.Contains(err.Error(), "NotFound") {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	movieDetails := movieModel.MovieDetails{
		Rating:   rating,
		Metadata: *metadata,
	}

	return &movieDetails, nil
}
