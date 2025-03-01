package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCalculateUpdatedRating(t *testing.T) {
	aggregatedRating := CreateAggregatedRating(uuid.NewString(), "movie")
	tests := []struct {
		testName             string
		recordID             string
		recordType           string
		ratingValue          int
		expectedRatingValue  float64
		expectedAmountRating int
	}{
		{"Test_First_Rating_Submission", aggregatedRating.RecordID, aggregatedRating.RecordType, 4, 4.0, 1},
		{"Test_Increase_Rating_Value", aggregatedRating.RecordID, aggregatedRating.RecordType, 7, 5.5, 2},
		{"Test_Decrease_Rating_Value", aggregatedRating.RecordID, aggregatedRating.RecordType, 1, 4.0, 3},
	}

	for _, test := range tests {
		aggregatedRating.CalculateUpdatedRating(test.ratingValue)
		require.Equal(t, test.expectedRatingValue, aggregatedRating.Rating)
		require.Equal(t, test.expectedAmountRating, aggregatedRating.AmountRatings)
	}
}
