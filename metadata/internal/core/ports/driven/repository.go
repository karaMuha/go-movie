package driven

import (
	"context"

	metadataModel "github.com/karaMuha/go-movie/metadata/pkg"
)

type IMetadataRepository interface {
	Save(ctx context.Context, metadata *metadataModel.Metadata) (*metadataModel.Metadata, error)
	Load(ctx context.Context, id string) (*metadataModel.Metadata, error)
}
