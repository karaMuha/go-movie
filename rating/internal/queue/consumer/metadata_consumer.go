package consumer

import (
	"context"
	"encoding/json"
	"log"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/segmentio/kafka-go"
)

type MetadataEventConsumer struct {
	app                     driving.IApplication
	Reader                  *kafka.Reader
	metadataEventRepository driven.IMetadataEventRepository
}

func NewMetadataEventConsumer(
	app driving.IApplication,
	address string,
	topic string,
	groupID string,
	metadataEventRepository driven.IMetadataEventRepository,
) *MetadataEventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
	})
	return &MetadataEventConsumer{
		app:                     app,
		Reader:                  reader,
		metadataEventRepository: metadataEventRepository,
	}
}

func (c *MetadataEventConsumer) StartReadingMetadataEvents() {
	for {
		message, err := c.Reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("error reading message from queue: %v\n", err)
			continue
		}
		err = c.Reader.CommitMessages(context.Background(), message)
		if err != nil {
			log.Printf("Error commiting fetched message: %v\n", err)
		}

		var event metadataModel.MetadataEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Unmarshal error on event: %v", err)
			continue
		}

		log.Printf("Read message: %v\n", event)
		respErr := c.app.SubmitMetadata(&ratingmodel.AggregatedRating{
			ID:            event.ID,
			RecordType:    string(event.RecordType),
			Rating:        0.0,
			AmountRatings: 0,
		})
		if respErr != nil {
			respErr = c.metadataEventRepository.Save(context.Background(), event)
			if respErr != nil {
				log.Printf("event with ID %s failed to be processed: %v\n", event.ID, err)
			}
		}
	}
}
