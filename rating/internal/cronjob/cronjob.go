package cronjob

import (
	"context"
	"log"
	"time"

	"github.com/karaMuha/go-movie/rating/internal/core"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type Cronjob struct {
	metadataEventRepo driven.IMetadataEventRepository
	app               core.Application
}

func NewCronjob(metadataEventRepo driven.IMetadataEventRepository, app core.Application) Cronjob {
	return Cronjob{
		metadataEventRepo: metadataEventRepo,
		app:               app,
	}
}

func (c *Cronjob) RunMetadata() {
	for {
		events, err := c.metadataEventRepo.Load(context.Background())
		if err != nil {
			log.Println(err)
		}

		for _, event := range events {
			err = c.app.SubmitMetadata(&ratingmodel.AggregatedRating{
				ID:            event.ID,
				RecordType:    string(event.RecordType),
				Rating:        0.0,
				AmountRatings: 0,
			})
			if err != nil {
				continue
			}
			err = c.metadataEventRepo.Delete(context.Background(), event.ID, string(event.RecordType))
			if err != nil {
				log.Printf("Saved event with ID %s and record_type %s processed but could not be cleaned up: %v\n", event.ID, event.RecordType, err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
