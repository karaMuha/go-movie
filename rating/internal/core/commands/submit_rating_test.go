package commands

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/karaMuha/go-movie/pkg/database/postgres"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	postgres_repo "github.com/karaMuha/go-movie/rating/internal/repository/postgres"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubmitRatingTestSuite struct {
	suite.Suite
	ctx                  context.Context
	cmd                  SubmitRatingCommand
	dbHandler            *sql.DB
	ratingRepo           driven.IRatingRepository
	aggregatedRatingRepo driven.IAggregatedRatingRepository
}

func TestSubmitRatingSuite(t *testing.T) {
	suite.Run(t, &SubmitRatingTestSuite{})
}

func (s *SubmitRatingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	initScriptPath := filepath.Join("..", "..", "..", "..", "dbscripts", "ratings_public_schema.sql")
	dbHandler, err := postgres.CreateTestContainer(s.ctx, "ratings_db", initScriptPath)
	require.NoError(s.T(), err)
	s.dbHandler = dbHandler
	ratingsRepository := postgres_repo.NewRatingRepository(s.dbHandler)
	aggregatedRatingRepo := postgres_repo.NewAggregatedRatingRepository(s.dbHandler)
	s.cmd = NewSubmitRatingCommand(&ratingsRepository, &aggregatedRatingRepo)
	s.ratingRepo = &ratingsRepository
	s.aggregatedRatingRepo = &aggregatedRatingRepo
}

// clear tables between tests to avoid conflicts and side effects
func (s *SubmitRatingTestSuite) AfterTest(suiteName, testName string) {
	queryClearRatingsTable := `DELETE FROM ratings`

	_, err := s.dbHandler.ExecContext(s.ctx, queryClearRatingsTable)
	require.NoError(s.T(), err)
}

func (s *SubmitRatingTestSuite) TestSubmitRating() {
	aggregatedRating := ratingmodel.AggregatedRating{
		ID:            uuid.NewString(),
		RecordType:    "movie",
		Rating:        7.0,
		AmountRatings: 1,
	}
	s.aggregatedRatingRepo.Save(s.ctx, &aggregatedRating)

	rating := ratingmodel.Rating{
		RecordID:   aggregatedRating.ID,
		RecordType: "movie",
		UserID:     "123",
		Value:      5,
	}
	respErr := s.cmd.SubmitRating(s.ctx, ratingmodel.RecordID(rating.RecordID), ratingmodel.RecordType(rating.RecordType), &rating)
	require.Nil(s.T(), respErr)

	updatedAggregatedRating, respErr := s.aggregatedRatingRepo.Load(s.ctx, aggregatedRating.ID, aggregatedRating.RecordType)
	require.Nil(s.T(), respErr)
	require.Equal(s.T(), 6.0, updatedAggregatedRating.Rating)
	require.Equal(s.T(), 2, updatedAggregatedRating.AmountRatings)
}
