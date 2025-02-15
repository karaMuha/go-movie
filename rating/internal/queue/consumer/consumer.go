package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/segmentio/kafka-go"
)

type MessageConsumer struct {
	app    driving.IApplication
	Reader *kafka.Reader
}

func NewMessageConsumer(app driving.IApplication, address string, topic string, gourpID string) *MessageConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
		GroupID: gourpID,
	})
	return &MessageConsumer{
		app:    app,
		Reader: reader,
	}
}

func (c *MessageConsumer) StartReading() {
	for {
		message, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message from queue: %v\n", err)
			continue
		}

		var event ratingmodel.RatingEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Unmarshal error on event: %v", err)
			continue
		}

		log.Printf("Read message: %v\n", event)

		c.app.SubmitRating(context.Background(), ratingmodel.RecordID(event.RecordID), ratingmodel.RecordType(event.RecordType), &ratingmodel.Rating{
			RecordID:   event.RecordID,
			RecordType: event.RecordType,
			UserID:     event.UserID,
			Value:      event.Value,
		})
	}
}
