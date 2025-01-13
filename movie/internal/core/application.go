package core

import (
	"github.com/karaMuha/go-movie/movie/internal/core/commands"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/movie/internal/core/queries"
)

type Application struct {
	appCommands
	appQueries
}

var _ driving.IApplication = (*Application)(nil)

type appCommands struct {
	commands.SubmitRatingCommand
}

type appQueries struct {
	queries.GetMovieDetailsQuery
}

func New(
	metadataGateway driven.IMetadataGateway,
	ratingGateway driven.IRatingGateway,
	messageProducer driven.IMessageProducer,
) Application {
	return Application{
		appCommands: appCommands{
			SubmitRatingCommand: commands.NewSubmitRatingCommand(ratingGateway, messageProducer),
		},
		appQueries: appQueries{
			GetMovieDetailsQuery: queries.NewGetMovieDetailsQuery(metadataGateway, ratingGateway),
		},
	}
}
