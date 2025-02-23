package driven

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
	ratingmodel "github.com/karaMuha/go-movie/rating/pkg"
)

type IMetadataGateway interface {
	GetMetadata(ctx context.Context, movieID string) (*metadataModel.Metadata, *dtos.RespErr)
	SubmitMetadata(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr)
}

type IRatingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, int, *dtos.RespErr)
	SubmitRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) *dtos.RespErr
}
