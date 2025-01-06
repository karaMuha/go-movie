package driven

import (
	"context"

	model "github.com/karaMuha/go-movie/metadata/metadataModel"
)

type IMetadataRepository interface {
	Save(ctx context.Context, id string, metadata *model.Metadata) error
	Load(ctx context.Context, id string) (*model.Metadata, error)
}
