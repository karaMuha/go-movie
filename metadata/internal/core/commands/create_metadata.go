package commands

import (
	"context"
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
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

func (c *CraeteMetadataCommand) CreateMetadata(ctx context.Context, cmd *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	err := domain.CreateMetadata(cmd.Title, cmd.Director)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: err.Error(),
		}
	}

	metadata, respErr := c.metadataRepo.Save(ctx, cmd)
	if respErr != nil {
		return nil, respErr
	}

	event := metadataModel.MetadataEvent{
		ID:         metadata.ID,
		RecordType: metadataModel.RecordTypeMovie,
		EventType:  metadataModel.MetadataEventTypeSubmitted,
	}
	err = c.producer.PublishMetadataSubmittedEvent(event)
	if err != nil {
		respErr = c.metadataEventRepo.Save(ctx, &event)
		if respErr != nil {
			log.Println(respErr)
		}
	}

	return metadata, nil
}
