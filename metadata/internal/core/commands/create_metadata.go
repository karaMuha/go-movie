package commands

import (
	"context"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type CraeteMetadataCommand struct {
	repo     driven.IMetadataRepository
	producer driven.IMessageProducer
}

func NewCreateMetadataCommand(repo driven.IMetadataRepository, producer driven.IMessageProducer) CraeteMetadataCommand {
	return CraeteMetadataCommand{
		repo:     repo,
		producer: producer,
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

	event := metadataModel.MetadataEvent{
		ID:         metadata.ID,
		RecordType: metadataModel.RecordTypeMovie,
		EventType:  metadataModel.MetadataEventTypeSubmitted,
	}
	_ = c.producer.PublishMetadataSubmittedEvent(event)

	return metadata, nil
}
