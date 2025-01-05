package core

import (
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/metadata/internal/core/queries"
)

type Application struct {
	appCommands
	appQueries
}

type appCommands struct{}

type appQueries struct {
	queries.GetMetadataQuery
}

var _ driving.IApplication = (*Application)(nil)

func New(repo driven.IRepository) Application {
	return Application{
		appQueries: appQueries{
			GetMetadataQuery: queries.NewGetMetadataQuery(repo),
		},
	}
}
