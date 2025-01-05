package queries

import (
	"context"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type GetMetadataQuery struct {
	repo driven.IRepository
}

func NewGetMetadataQuery(repo driven.IRepository) GetMetadataQuery {
	return GetMetadataQuery{
		repo: repo,
	}
}

func (q *GetMetadataQuery) GetMetadata(ctx context.Context, id string) (*model.Metadata, error) {
	return q.repo.Load(ctx, id)
}
