package cronjob

import (
	"context"
	"log"
	"time"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
)

type Cronjob struct {
	metadataEventRepo driven.IMetadataEventRepository
	producer          driven.IMessageProducer
	doneChan          chan struct{}
}

func NewCronjob(metadataEventRepo driven.IMetadataEventRepository, producer driven.IMessageProducer) Cronjob {
	return Cronjob{
		metadataEventRepo: metadataEventRepo,
		producer:          producer,
	}
}

func (c *Cronjob) Run() {
	for {
		select {
		case <-c.doneChan:
			log.Println("Stopping cronjob gracefully")
			return
		default:
			time.Sleep(1 * time.Minute)

			events, respErr := c.metadataEventRepo.Load(context.Background())
			if respErr != nil {
				log.Println(respErr)
				continue
			}

			for _, event := range events {
				err := c.producer.PublishMetadataSubmittedEvent(event)
				if respErr != nil {
					log.Println(err)
					continue
				}
				respErr = c.metadataEventRepo.Delete(context.Background(), event.ID)
				if respErr != nil {
					log.Printf("Saved event with ID %s and record type %s published but could not be cleaned up: %v\n", event.ID, event.RecordType, respErr)
				}
			}
		}
	}
}

func (c *Cronjob) GracefulStop() {
	c.doneChan <- struct{}{}
	close(c.doneChan)
}
