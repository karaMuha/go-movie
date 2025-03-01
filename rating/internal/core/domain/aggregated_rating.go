package domain

type AggregatedRating struct {
	RecordID      string
	RecordType    string
	Rating        float64
	AmountRatings int
}

func CreateAggregatedRating(recordID, recordType string) *AggregatedRating {
	return &AggregatedRating{
		RecordID:      recordID,
		RecordType:    recordType,
		Rating:        0.0,
		AmountRatings: 0,
	}
}

func (r *AggregatedRating) CalculateUpdatedRating(ratingValue float64) {
	ratingSum := r.Rating * float64(r.AmountRatings)
	ratingSum += float64(ratingValue)
	newRating := ratingSum / (float64(r.AmountRatings) + 1.0)

	r.Rating = newRating
	r.AmountRatings += 1
}
