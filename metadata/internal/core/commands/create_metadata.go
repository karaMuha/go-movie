package commands

import (
	"context"
	"log"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type CraeteMetadataCommand struct {
	metadataRepo      driven.IMetadataRepository
	metadataEventRepo driven.IMetadataEventRepository
	producer          driven.IMessageProducer
}

func NewCreateMetadataCommand(
	repo driven.IMetadataRepository,
	producer driven.IMessageProducer,
	metadataEventRepo driven.IMetadataEventRepository,
) CraeteMetadataCommand {
	return CraeteMetadataCommand{
		metadataRepo:      repo,
		metadataEventRepo: metadataEventRepo,
		producer:          producer,
	}
}

func (c *CraeteMetadataCommand) CreateMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, error) {
	err := domain.CreateMetadata(cmd.Title, cmd.Director)
	if err != nil {
		return nil, err
	}

	metadata, err := c.metadataRepo.Save(ctx, cmd)
	if err != nil {
		return nil, err
	}

	event := metadataModel.MetadataEvent{
		ID:         metadata.ID,
		RecordType: metadataModel.RecordTypeMovie,
		EventType:  metadataModel.MetadataEventTypeSubmitted,
	}
	err = c.producer.PublishMetadataSubmittedEvent(event)
	if err != nil {
		err = c.metadataEventRepo.Save(ctx, &event)
		if err != nil {
			log.Println(err)
		}
	}

	return metadata, nil
}
