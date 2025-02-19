package core

import (
	"github.com/karaMuha/go-movie/metadata/internal/core/commands"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/metadata/internal/core/queries"
)

type Application struct {
	appCommands
	appQueries
}

type appCommands struct {
	commands.CraeteMetadataCommand
}

type appQueries struct {
	queries.GetMetadataQuery
}

var _ driving.IApplication = (*Application)(nil)

func New(repo driven.IMetadataRepository, producer driven.IMessageProducer) Application {
	return Application{
		appCommands: appCommands{
			CraeteMetadataCommand: commands.NewCreateMetadataCommand(repo, producer),
		},
		appQueries: appQueries{
			GetMetadataQuery: queries.NewGetMetadataQuery(repo),
		},
	}
}
