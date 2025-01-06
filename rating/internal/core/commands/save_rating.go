package commands

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/ratingModel"
)

type SaveRatingCommand struct {
	ratingsRepo driven.IRatingRepository
}

func NewSaveRatingCommand(ratingsRepo driven.IRatingRepository) SaveRatingCommand {
	return SaveRatingCommand{
		ratingsRepo: ratingsRepo,
	}
}

func (c *SaveRatingCommand) SaveRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.ratingsRepo.Save(ctx, recordID, recordType, rating)
}
