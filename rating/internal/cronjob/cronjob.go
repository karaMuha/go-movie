package cronjob

import (
	"context"
	"log"
	"time"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type Cronjob struct {
	metadataEventRepo driven.IMetadataEventRepository
	app               driving.IApplication
	doneChan          chan struct{}
}

func NewCronjob(metadataEventRepo driven.IMetadataEventRepository, app driving.IApplication) Cronjob {
	return Cronjob{
		metadataEventRepo: metadataEventRepo,
		app:               app,
		doneChan:          make(chan struct{}),
	}
}

func (c *Cronjob) RunMetadata() {
	for {
		select {
		case <-c.doneChan:
			log.Println("Stopping cronjob gracefully")
			return
		default:
			time.Sleep(10 * time.Second)

			events, err := c.metadataEventRepo.Load(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}

			for _, event := range events {
				err = c.app.SubmitMetadata(&ratingmodel.AggregatedRating{
					ID:            event.ID,
					RecordType:    string(event.RecordType),
					Rating:        0.0,
					AmountRatings: 0,
				})
				if err != nil {
					log.Println(err)
					continue
				}
				err = c.metadataEventRepo.Delete(context.Background(), event.ID, string(event.RecordType))
				if err != nil {
					log.Printf("Saved event with ID %s and record_type %s processed but could not be cleaned up: %v\n", event.ID, event.RecordType, err)
				}
			}
		}
	}
}

func (c *Cronjob) GracefulStop() {
	c.doneChan <- struct{}{}
}
