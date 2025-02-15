package driving

import (
	"context"

	model "github.com/karaMuha/go-movie/metadata/pkg"
)

type IApplication interface {
	ICommands
	IQueries
}

type ICommands interface{}

type IQueries interface {
	GetMetadata(ctx context.Context, id string) (*model.Metadata, error)
}
