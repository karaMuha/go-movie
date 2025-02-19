package producer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/segmentio/kafka-go"
)

type MessageProducer struct {
	Writer *kafka.Writer
}

var _ driven.IMessageProducer = (*MessageProducer)(nil)

func NewMessageProducer(address string, topic string) *MessageProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{address},
		Topic:   topic,
	})
	writer.AllowAutoTopicCreation = true
	writer.BatchSize = 1

	return &MessageProducer{
		Writer: writer,
	}
}

func (p *MessageProducer) PublishMetadataSubmittedEvent(event metadataModel.MetadataEvent) error {
	encodedEvent, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error encoding event before publishing: %v\n", err)
		return nil
	}

	message := kafka.Message{
		Value: []byte(encodedEvent),
	}

	err = p.Writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Error publishing event: %v\n", err)
		return err
	}

	return nil
}
