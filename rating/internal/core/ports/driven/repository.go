package driven

import (
	"context"

	model "github.com/karaMuha/go-movie/rating/pkg"
)

type IRatingRepository interface {
	Load(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error)
	Save(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}
