package postgres_repo

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
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

func (r *MetadataEventRepository) Save(ctx context.Context, event *metadataModel.MetadataEvent) *dtos.RespErr {
	query := `
		INSERT INTO metadata_created_events(id, record_type, event_type)
		VALUES($1, $2, $3);
	`
	_, err := r.db.ExecContext(ctx, query, event.ID, event.RecordType, event.EventType)
	if err != nil {
		log.Println(err)
		return &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	return nil
}

func (r *MetadataEventRepository) Load(ctx context.Context) ([]metadataModel.MetadataEvent, *dtos.RespErr) {
	query := `
		SELECT *
		FROM metadata_created_events
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: err.Error(),
		}
	}

	var events []metadataModel.MetadataEvent
	for rows.Next() {
		var event metadataModel.MetadataEvent
		if err := rows.Scan(
			&event.ID,
			&event.RecordType,
			&event.EventType,
		); err != nil {
			return nil, &dtos.RespErr{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: err.Error(),
			}
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *MetadataEventRepository) Delete(ctx context.Context, ID string) *dtos.RespErr {
	query := `
		DELETE FROM metadata_created_events
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, ID)
	return &dtos.RespErr{
		StatusCode:    http.StatusInternalServerError,
		StatusMessage: err.Error(),
	}
}
