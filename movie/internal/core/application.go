package core

import (
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/movie/internal/core/queries"
)

type Application struct {
	appCommands
	appQueries
}

var _ driving.IApplication = (*Application)(nil)

type appCommands struct{}

type appQueries struct {
	queries.GetMovieDetailsQuery
}

func New(metadataGateway driven.IMetadataGateway, ratingGateway driven.IRatingGateway) Application {
	return Application{
		appCommands: appCommands{},
		appQueries: appQueries{
			GetMovieDetailsQuery: queries.NewGetMovieDetailsQuery(metadataGateway, ratingGateway),
		},
	}
}
