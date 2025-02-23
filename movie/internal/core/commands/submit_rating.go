package commands

import (
	"context"
	"net/http"

	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/dtos"
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

func (c *SubmitRatingCommand) SubmitRating(ctx context.Context, cmd *ratingmodel.Rating) *dtos.RespErr {
	/* err := c.ratingGateway.SubmitRating(ctx, ratingmodel.RecordID(cmd.RecordID), ratingmodel.RecordType(cmd.RecordType), cmd)
	if err != nil {
		return err
	} */

	event := &ratingmodel.RatingEvent{
		RecordID:   cmd.RecordID,
		RecordType: cmd.RecordType,
		UserID:     cmd.UserID,
		Value:      cmd.Value,
		EventType:  ratingmodel.RatingEventTypeSubmit,
	}

	err := c.messageProducer.PublishRatingSubmittedEvent(*event)
	if err != nil {
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	return nil
}
