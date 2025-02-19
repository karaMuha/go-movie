package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
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

func (r *MetadataRepository) Load(ctx context.Context, id string) (*metadataModel.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return m, nil
}

func (r *MetadataRepository) Save(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, error) {
	id := uuid.NewString()
	metadata.ID = id
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return metadata, nil
}
