package commands

import (
	"context"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type CraeteMetadataCommand struct {
	repo driven.IMetadataRepository
}

func NewCreateMetadataCommand(repo driven.IMetadataRepository) CraeteMetadataCommand {
	return CraeteMetadataCommand{
		repo: repo,
	}
}

func (c *CraeteMetadataCommand) CreateMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, error) {
	err := domain.CreateMetadata(cmd.Title, cmd.Director)
	if err != nil {
		return nil, err
	}

	metadata, err := c.repo.Save(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
