package driving

import (
	"context"

	model "github.com/karaMuha/go-movie/movie/movieModel"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface{}

type IQueries interface {
	GetMovieDetails(ctx context.Context, movieID string) (*model.MovieDetails, error)
}
