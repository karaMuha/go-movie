package commands

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitMetadataCommand struct {
	metadataRepo driven.IMetadatarepository
}

func NewSubmitMetadataCommand(metadataRepo driven.IMetadatarepository) SubmitMetadataCommand {
	return SubmitMetadataCommand{
		metadataRepo: metadataRepo,
	}
}

func (c *SubmitMetadataCommand) SubmitMetadata(cmd *ratingmodel.AggregatedRating) {
	_ = c.metadataRepo.Save(context.Background(), cmd)
}
