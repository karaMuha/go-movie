package driven

import (
	"context"

	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type IRepository interface {
	Save(ctx context.Context, id string, metadata *model.Metadata) error
	Load(ctx context.Context, id string) (*model.Metadata, error)
}
