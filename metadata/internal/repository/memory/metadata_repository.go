package memory

import (
	"context"
	"sync"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type MetadataRepository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

var _ driven.IMetadataRepository = (*MetadataRepository)(nil)

func New() *MetadataRepository {
	return &MetadataRepository{
		data: make(map[string]*model.Metadata),
	}
}

func (r *MetadataRepository) Load(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return m, nil
}

func (r *MetadataRepository) Save(ctx context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
