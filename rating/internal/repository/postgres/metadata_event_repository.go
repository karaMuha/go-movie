package postgres_repo

import (
	"context"
	"database/sql"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
)

type MetadataEventRepository struct {
	db *sql.DB
}

func NewMetadataEventRepository(db *sql.DB) MetadataEventRepository {
	return MetadataEventRepository{
		db: db,
	}
}

var _ driven.IMetadataEventRepository = (*MetadataEventRepository)(nil)

func (r *MetadataEventRepository) Save(ctx context.Context, event metadataModel.MetadataEvent) error {
	query := `
		INSERT INTO metadata_created_event (id, record_type, event_type)
		VALUES ($1, $2, $3);
	`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.RecordType, event.EventType)
	return err
}

func (r *MetadataEventRepository) Load(ctx context.Context) ([]metadataModel.MetadataEvent, error) {
	query := `
		SELECT *
		FROM metadata_created_event
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var events []metadataModel.MetadataEvent
	for rows.Next() {
		var event metadataModel.MetadataEvent
		if err := rows.Scan(
			&event.ID,
			&event.RecordType,
			&event.EventType,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *MetadataEventRepository) Delete(ctx context.Context, ID, recordType string) error {
	query := `
		DELETE FROM metadata_created_events
		WHERE id = $1 AND record_type = $2
	`
	_, err := r.db.ExecContext(ctx, query, ID, recordType)
	return err
}
