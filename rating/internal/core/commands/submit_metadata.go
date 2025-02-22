package commands

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitMetadataCommand struct {
	metadataRepo driven.IAggregatedRatingRepository
}

func NewSubmitMetadataCommand(metadataRepo driven.IAggregatedRatingRepository) SubmitMetadataCommand {
	return SubmitMetadataCommand{
		metadataRepo: metadataRepo,
	}
}

func (c *SubmitMetadataCommand) SubmitMetadata(cmd *ratingmodel.AggregatedRating) error {
	return c.metadataRepo.Save(context.Background(), cmd)
}
