package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/segmentio/kafka-go"
)

type QueueConsumer struct {
	app    driving.IApplication
	reader *kafka.Reader
}

func NewQueueConsumer(app driving.IApplication, address string, topic string, gourpID string) *QueueConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
		GroupID: gourpID,
	})
	return &QueueConsumer{
		app:    app,
		reader: reader,
	}
}

func (c *QueueConsumer) StartReading() {
	message, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		log.Printf("error reading message from queue: %v\n", err)
	}

	var event ratingmodel.RatingEvent
	err = json.Unmarshal(message.Value, &event)
	if err != nil {
		log.Printf("Unmarshal error on event: %v", err)
	}

	c.app.SubmitRating(context.Background(), ratingmodel.RecordID(event.RecordID), ratingmodel.RecordType(event.RecordType), &ratingmodel.Rating{
		RecordID:   event.RecordID,
		RecordType: event.RecordType,
		UserID:     event.UserID,
		Value:      event.Value,
	})
}
