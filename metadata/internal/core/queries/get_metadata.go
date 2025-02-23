package queries

import (
	"context"

	"github.com/karaMuha/go-movie/metadata/internal/core/ports/driven"
	model "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type GetMetadataQuery struct {
	metadataRepo driven.IMetadataRepository
}

func NewGetMetadataQuery(metadataRepo driven.IMetadataRepository) GetMetadataQuery {
	return GetMetadataQuery{
		metadataRepo: metadataRepo,
	}
}

func (q *GetMetadataQuery) GetMetadata(ctx context.Context, id string) (*model.Metadata, *dtos.RespErr) {
	return q.metadataRepo.Load(ctx, id)
}
