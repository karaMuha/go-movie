package queries

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/movie/internal/core/ports/driven"
)

type GetMetadataQuery struct {
	metadataGateway driven.IMetadataGateway
}

func NewGetMetadataQuery(metadataGateway driven.IMetadataGateway) GetMetadataQuery {
	return GetMetadataQuery{
		metadataGateway: metadataGateway,
	}
}

func (q *GetMetadataQuery) GetMetadata(ctx context.Context, ID string) (*metadataModel.Metadata, error) {
	return q.metadataGateway.GetMetadata(ctx, ID)
}
