package postgres_repo

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
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

func (m *MetadataRepository) Load(ctx context.Context, id string) (*metadataModel.Metadata, *dtos.RespErr) {
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
			return nil, &dtos.RespErr{
				StatusCode:    http.StatusNotFound,
				StatusMessage: "Not Found",
			}
		}
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	return &metadata, nil
}

func (m *MetadataRepository) Save(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	query := `
		INSERT INTO metadata (title, description, director)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	row := m.db.QueryRowContext(ctx, query, metadata.Title, metadata.Description, metadata.Director)

	var id string
	if err := row.Scan(&id); err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	metadata.ID = id

	return metadata, nil
}
