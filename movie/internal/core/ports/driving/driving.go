package driving

import (
	"context"

	model "github.com/karaMuha/go-movie/movie/movieModel"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SubmitRating(ctx context.Context, cmd *ratingmodel.Rating) error
}

type IQueries interface {
	GetMovieDetails(ctx context.Context, movieID string) (*model.MovieDetails, error)
}
