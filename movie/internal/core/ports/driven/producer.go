package driven

import ratingmodel "github.com/karaMuha/go-movie/rating/pkg"

type IMessageProducer interface {
	PublishRatingSubmittedEvent(event ratingmodel.RatingEvent) error
}
