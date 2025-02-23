package commands

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-movie/pkg/dtos"
	"github.com/karaMuha/go-movie/rating/internal/core/domain"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitRatingCommand struct {
	ratingsRepo          driven.IRatingRepository
	aggregatedRatingRepo driven.IAggregatedRatingRepository
}

func NewSubmitRatingCommand(ratingsRepo driven.IRatingRepository, metadataRepo driven.IAggregatedRatingRepository) SubmitRatingCommand {
	return SubmitRatingCommand{
		ratingsRepo:          ratingsRepo,
		aggregatedRatingRepo: metadataRepo,
	}
}

func (c *SubmitRatingCommand) SubmitRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) *dtos.RespErr {
	err := domain.SubmitRating(rating.RecordID, rating.RecordType, rating.UserID, rating.Value)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: err.Error(),
		}
	}

	respErr := c.ratingsRepo.Save(ctx, recordID, recordType, rating)
	if respErr != nil {
		return respErr
	}

	aggregatedRating, respErr := c.aggregatedRatingRepo.Load(ctx, string(recordID), string(recordType))
	if respErr != nil {
		// save in table for cronjob
		return respErr
	}

	ratingSum := aggregatedRating.Rating * float64(aggregatedRating.AmountRatings)
	ratingSum += float64(rating.Value)
	newRating := ratingSum / (float64(aggregatedRating.AmountRatings) + 1.0)
	aggregatedRating.AmountRatings += 1
	aggregatedRating.Rating = newRating
	respErr = c.aggregatedRatingRepo.Update(ctx, aggregatedRating)
	if respErr != nil {
		// save in table for cronjob
	}

	return nil
}
