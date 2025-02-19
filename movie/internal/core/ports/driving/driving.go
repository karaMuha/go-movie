package driving

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	model "github.com/karaMuha/go-movie/movie/movieModel"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SubmitRating(ctx context.Context, cmd *ratingmodel.Rating) error
	SubmitMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, error)
}

type IQueries interface {
	GetMovieDetails(ctx context.Context, movieID string) (*model.MovieDetails, error)
	GetMetadata(ctx context.Context, ID string) (*metadataModel.Metadata, error)
}
