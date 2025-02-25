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
	commands.CreateMetadataCommand
}

type appQueries struct {
	queries.GetMetadataQuery
}

var _ driving.IApplication = (*Application)(nil)

func New(metadataRepo driven.IMetadataRepository,
	producer driven.IMessageProducer,
	metadataEventRepo driven.IMetadataEventRepository,
) Application {
	return Application{
		appCommands: appCommands{
			CreateMetadataCommand: commands.NewCreateMetadataCommand(metadataRepo, producer, metadataEventRepo),
		},
		appQueries: appQueries{
			GetMetadataQuery: queries.NewGetMetadataQuery(metadataRepo),
		},
	}
}
