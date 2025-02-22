package commands

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/karaMuha/go-movie/pkg/database/postgres"
	postgres_repo "github.com/karaMuha/go-movie/rating/internal/repository/postgres"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubmitRatingTestSuite struct {
	suite.Suite
	ctx       context.Context
	cmd       SubmitRatingCommand
	dbHandler *sql.DB
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
	metadataRepository := postgres_repo.NewAggregatedMetadataRepository(s.dbHandler)
	s.cmd = NewSubmitRatingCommand(&ratingsRepository, &metadataRepository)
}

// clear tables between tests to avoid conflicts and side effects
func (s *SubmitRatingTestSuite) AfterTest(suiteName, testName string) {
	queryClearRatingsTable := `DELETE FROM ratings`

	_, err := s.dbHandler.ExecContext(s.ctx, queryClearRatingsTable)
	require.NoError(s.T(), err)
}

func (s *SubmitRatingTestSuite) TestSubmitRating() {
	rating := ratingmodel.Rating{
		RecordID:   "123",
		RecordType: "movie",
		UserID:     "123",
		Value:      5,
	}
	err := s.cmd.SubmitRating(s.ctx, "123", "movie", &rating)
	require.NoError(s.T(), err)
}
