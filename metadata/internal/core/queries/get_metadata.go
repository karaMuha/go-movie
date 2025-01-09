package queries

import (
	"context"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type GetMetadataQuery struct {
	metadataRepo driven.IMetadataRepository
}

func NewGetMetadataQuery(metadataRepo driven.IMetadataRepository) GetMetadataQuery {
	return GetMetadataQuery{
		metadataRepo: metadataRepo,
	}
}

func (q *GetMetadataQuery) GetMetadata(ctx context.Context, id string) (*model.Metadata, error) {
	return q.metadataRepo.Load(ctx, id)
}
