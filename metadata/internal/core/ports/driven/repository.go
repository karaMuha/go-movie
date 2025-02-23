package driven

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

type IMetadataRepository interface {
	Save(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, *dtos.RespErr)
	Load(ctx context.Context, id string) (*metadataModel.Metadata, *dtos.RespErr)
}

type IMetadataEventRepository interface {
	Save(ctx context.Context, event *metadataModel.MetadataEvent) *dtos.RespErr
	Load(ctx context.Context) (*[]metadataModel.MetadataEvent, *dtos.RespErr)
	Delete(ctx context.Context, ID string) *dtos.RespErr
}
