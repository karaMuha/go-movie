package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/segmentio/kafka-go"
)

type RatingEventConsumer struct {
	app          driving.IApplication
	RatingReader *kafka.Reader
	doneChan     chan struct{}
}

func NewRatingEventConsumer(app driving.IApplication, address string, topic string, gourpID string) *RatingEventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
	})
	return &RatingEventConsumer{
		app:          app,
		RatingReader: reader,
	}
}

func (c *RatingEventConsumer) StartReadingRatingEvents() {
	for {
		select {
		case <-c.doneChan:
			log.Println("Stopped reading rating events")
			return
		default:
			message, err := c.RatingReader.FetchMessage(context.Background())
			if err != nil {
				log.Printf("error reading rating message from queue: %v\n", err)
				continue
			}

			go c.ProcessEvent(message)
		}
	}
}

func (c *RatingEventConsumer) StopReadingEvents() {
	c.doneChan <- struct{}{}
}

func (c *RatingEventConsumer) ProcessEvent(message kafka.Message) {
	err := c.RatingReader.CommitMessages(context.Background(), message)
	if err != nil {
		log.Printf("Error commiting fetched message: %v\n", err)
		return
	}

	var event ratingmodel.RatingEvent
	err = json.Unmarshal(message.Value, &event)
	if err != nil {
		log.Printf("Unmarshal error on event: %v", err)
		return
	}

	log.Printf("Read message: %v\n", event)

	rating := &ratingmodel.Rating{
		RecordID:   event.RecordID,
		RecordType: event.RecordType,
		UserID:     event.UserID,
		Value:      event.Value,
	}

	c.app.SubmitRating(context.Background(), ratingmodel.RecordID(event.RecordID), ratingmodel.RecordType(event.RecordType), rating)
}
