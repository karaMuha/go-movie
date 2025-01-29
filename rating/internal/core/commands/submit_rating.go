package commands

import (
	"context"

	"github.com/karaMuha/go-movie/rating/internal/core/domain"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitRatingCommand struct {
	ratingsRepo driven.IRatingRepository
}

func NewSubmitRatingCommand(ratingsRepo driven.IRatingRepository) SubmitRatingCommand {
	return SubmitRatingCommand{
		ratingsRepo: ratingsRepo,
	}
}

func (c *SubmitRatingCommand) SubmitRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	err := domain.SubmitRating(rating.RecordID, rating.RecordType, rating.UserID, rating.Value)
	if err != nil {
		return err
	}
	return c.ratingsRepo.Save(ctx, recordID, recordType, rating)
}
