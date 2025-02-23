package driving

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface {
	CreateMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr)
}

type IQueries interface {
	GetMetadata(ctx context.Context, id string) (*metadataModel.Metadata, *dtos.RespErr)
}
