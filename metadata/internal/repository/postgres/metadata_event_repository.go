package postgres_repo

import (
	"context"
	"database/sql"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
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

func (r *MetadataEventRepository) Save(ctx context.Context, event *metadataModel.MetadataEvent) error {
	query := `
		INSERT INTO metadata_creted_events(id, record_ype, event_type)
		VALUES($1, $2, $3);
	`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.RecordType, event.EventType)
	if err != nil {
		return err
	}

	return nil
}

func (r *MetadataEventRepository) Load(ctx context.Context) (*[]metadataModel.MetadataEvent, error) {
	query := `
		SELECT *
		FROM metadata_created_events
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
	return &events, nil
}

func (r *MetadataEventRepository) Delete(ctx context.Context, ID string) error {
	query := `
		DELETE FROM metadata_created_events
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, ID)
	return err
}
