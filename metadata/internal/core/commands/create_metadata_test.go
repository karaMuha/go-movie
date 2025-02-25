package commands

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	postgres_repo "github.com/karaMuha/go-movie/metadata/internal/repository/postgres"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/database/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MessagePublisherMock struct{}

func (p *MessagePublisherMock) PublishMetadataSubmittedEvent(event metadataModel.MetadataEvent) error {
	if event.RecordType == "fail" {
		return errors.New("simulating failure")
	}
	return nil
}

type CreateMetadataTestSuite struct {
	suite.Suite
	ctx               context.Context
	cmd               CreateMetadataCommand
	db                *sql.DB
	metadataRepo      driven.IMetadataRepository
	producer          driven.IMessageProducer
	metadataEventRepo driven.IMetadataEventRepository
}

func TestCreateMetadataTestSuite(t *testing.T) {
	suite.Run(t, &CreateMetadataTestSuite{})
}

func (s *CreateMetadataTestSuite) SetupSuite() {
	s.ctx = context.Background()
	initScriptPath := filepath.Join("..", "..", "..", "..", "dbscripts", "metadata_public_schema.sql")
	db, err := postgres.CreateTestContainer(s.ctx, "metadata_db", initScriptPath)
	require.NoError(s.T(), err)
	s.db = db
	metadataRepo := postgres_repo.NewMetadataRepository(s.db)
	metadataEventRepo := postgres_repo.NewMetadataEventRepository(s.db)
	publisherMock := &MessagePublisherMock{}
	s.metadataRepo = &metadataRepo
	s.producer = publisherMock
	s.metadataEventRepo = &metadataEventRepo
	cmd := NewCreateMetadataCommand(&metadataRepo, publisherMock, &metadataEventRepo)
	s.cmd = cmd
}

func (s *CreateMetadataTestSuite) AfterTest(suiteName, testName string) {
	query := `
		DELETE FROM metadata
	`
	_, err := s.db.ExecContext(s.ctx, query)
	require.NoError(s.T(), err)
}

func (s *CreateMetadataTestSuite) TestCreateMetadataSuccess() {
	metadata := metadataModel.Metadata{
		Title:       "Batman",
		Description: "Good movie",
		Director:    "Christopher Nolan",
		RecordType:  "movie",
	}
	result, respErr := s.cmd.CreateMetadata(s.ctx, &metadata)
	require.Nil(s.T(), respErr)
	require.NotEmpty(s.T(), result.ID)
}

func (s *CreateMetadataTestSuite) TestCreateMetadataFailPublishingMessage() {
	metadata := metadataModel.Metadata{
		Title:       "Batman",
		Description: "Good movie",
		Director:    "Christopher Nolan",
		RecordType:  "fail",
	}
	result, respErr := s.cmd.CreateMetadata(s.ctx, &metadata)
	require.Nil(s.T(), respErr)
	require.NotEmpty(s.T(), result.ID)

	savedEvent, respErr := s.metadataEventRepo.Load(s.ctx)
	require.Nil(s.T(), respErr)
	require.Len(s.T(), savedEvent, 1)
	require.Equal(s.T(), result.ID, savedEvent[0].ID)
}
