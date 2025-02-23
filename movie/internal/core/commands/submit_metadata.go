package commands

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type SubmitMetadataCommand struct {
	metadataGateway driven.IMetadataGateway
}

func NewSubmitMetadataCommand(metadataGateway driven.IMetadataGateway) SubmitMetadataCommand {
	return SubmitMetadataCommand{
		metadataGateway: metadataGateway,
	}
}

func (c *SubmitMetadataCommand) SubmitMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	return c.metadataGateway.SubmitMetadata(ctx, cmd)
}
