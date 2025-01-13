package commands

import (
	"context"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type SubmitRatingCommand struct {
	ratingGateway   driven.IRatingGateway
	messageProducer driven.IMessageProducer
}

func NewSubmitRatingCommand(ratingGateway driven.IRatingGateway, messageProducer driven.IMessageProducer) SubmitRatingCommand {
	return SubmitRatingCommand{
		ratingGateway:   ratingGateway,
		messageProducer: messageProducer,
	}
}

func (c *SubmitRatingCommand) SubmitRating(ctx context.Context, cmd *ratingmodel.Rating) error {
	err := c.ratingGateway.SubmitRating(ctx, ratingmodel.RecordID(cmd.RecordID), ratingmodel.RecordType(cmd.RecordType), cmd)
	if err != nil {
		return err
	}

	event := &ratingmodel.RatingEvent{
		RecordID:   cmd.RecordID,
		RecordType: cmd.RecordType,
		UserID:     cmd.UserID,
		Value:      cmd.Value,
		EventType:  ratingmodel.RatingEventTypeSubmit,
	}

	c.messageProducer.PublishRatingSubmittedEvent(*event)

	return nil
}
