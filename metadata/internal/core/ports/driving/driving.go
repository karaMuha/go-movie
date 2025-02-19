package driving

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	CreateMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, error)
}

type IQueries interface {
	GetMetadata(ctx context.Context, id string) (*metadataModel.Metadata, error)
}
