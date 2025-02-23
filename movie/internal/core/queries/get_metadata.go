package queries

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type GetMetadataQuery struct {
	metadataGateway driven.IMetadataGateway
}

func NewGetMetadataQuery(metadataGateway driven.IMetadataGateway) GetMetadataQuery {
	return GetMetadataQuery{
		metadataGateway: metadataGateway,
	}
}

func (q *GetMetadataQuery) GetMetadata(ctx context.Context, ID string) (*metadataModel.Metadata, *dtos.RespErr) {
	return q.metadataGateway.GetMetadata(ctx, ID)
}
