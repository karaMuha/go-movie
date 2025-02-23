package memory

import (
	"context"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type MetadataRepository struct {
	sync.RWMutex
	data map[string]*metadataModel.Metadata
}

var _ driven.IMetadataRepository = (*MetadataRepository)(nil)

func New() *MetadataRepository {
	return &MetadataRepository{
		data: make(map[string]*metadataModel.Metadata),
	}
}

func (r *MetadataRepository) Load(ctx context.Context, id string) (*metadataModel.Metadata, *dtos.RespErr) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, &dtos.RespErr{
			StatusCode:    http.StatusNotFound,
			StatusMessage: "Not Found",
		}
	}

	return m, nil
}

func (r *MetadataRepository) Save(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr) {
	id := uuid.NewString()
	metadata.ID = id
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return metadata, nil
}
