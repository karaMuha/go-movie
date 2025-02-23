package driving

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	model "github.com/karaMuha/go-movie/movie/movieModel"
	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	SubmitRating(ctx context.Context, cmd *ratingmodel.Rating) *dtos.RespErr
	SubmitMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr)
}

type IQueries interface {
	GetMovieDetails(ctx context.Context, movieID string) (*model.MovieDetails, *dtos.RespErr)
	GetMetadata(ctx context.Context, ID string) (*metadataModel.Metadata, *dtos.RespErr)
}
