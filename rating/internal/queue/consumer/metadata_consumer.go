package consumer

import (
	"context"
	"encoding/json"
	"log"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/segmentio/kafka-go"
)

type MetadataEventConsumer struct {
	app    driving.IApplication
	Reader *kafka.Reader
}

func NewMetadataEventConsumer(app driving.IApplication, address string, topic string, gourpID string) *MetadataEventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
	})
	return &MetadataEventConsumer{
		app:    app,
		Reader: reader,
	}
}

func (c *MetadataEventConsumer) StartReadingMetadataEvents() {
	for {
		message, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message from queue: %v\n", err)
			continue
		}

		var event metadataModel.MetadataEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Unmarshal error on event: %v", err)
			continue
		}

		log.Printf("Read message: %v\n", event)
		c.app.SubmitMetadata(&ratingmodel.AggregatedRating{
			ID:            event.ID,
			RecordType:    string(event.RecordType),
			Rating:        0.0,
			AmountRatings: 0,
		})
	}
}
