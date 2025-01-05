package memory

import (
	"context"
	"sync"

	"github.com/karaMuha/go-movie/metadata/internal/core/domain"
	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

var _ driven.IRepository = (*Repository)(nil)

func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

func (r *Repository) Load(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return m, nil
}

func (r *Repository) Save(ctx context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
