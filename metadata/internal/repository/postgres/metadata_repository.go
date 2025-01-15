package postgres_repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type MetadataRepository struct {
	db *sql.DB
}

var _ driven.IMetadataRepository = (*MetadataRepository)(nil)

func NewMetadataRepository(db *sql.DB) MetadataRepository {
	return MetadataRepository{
		db: db,
	}
}

func (m *MetadataRepository) Load(ctx context.Context, id string) (*metadataModel.Metadata, error) {
	query := `
		SELECT *
		FROM metadata
		WHERE id = $1
	`

	var metadata metadataModel.Metadata
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&metadata.ID,
		&metadata.Title,
		&metadata.Description,
		&metadata.Director,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("metadata not found")
		}
		return nil, err
	}

	return &metadata, nil
}

func (m *MetadataRepository) Save(ctx context.Context, id string, metadata *metadataModel.Metadata) error {
	query := `
		INSERT INTO metadata (title, description, director)
		VALUES ($1, $2, $3)
	`
	_, err := m.db.ExecContext(ctx, query, metadata.Title, metadata.Description, metadata.Director)

	return err
}
