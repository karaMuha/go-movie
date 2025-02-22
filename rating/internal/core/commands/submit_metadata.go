package commands

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitMetadataCommand struct {
	aggregatedRatingRepo driven.IAggregatedRatingRepository
}

func NewSubmitMetadataCommand(metadataRepo driven.IAggregatedRatingRepository) SubmitMetadataCommand {
	return SubmitMetadataCommand{
		aggregatedRatingRepo: metadataRepo,
	}
}

func (c *SubmitMetadataCommand) SubmitMetadata(cmd *ratingmodel.AggregatedRating) error {
	return c.aggregatedRatingRepo.Save(context.Background(), cmd)
}
