package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) ListImages(
	ctx context.Context,
	ids []int,
	moreKeys ...string,
) ([]common.Image, error) {
	_, span := trace.StartSpan(ctx, "upload.storage.list")
	defer span.End()

	var result []common.Image

	if err := store.db.
		Table(common.Image{}.TableName()).
		Where("id in ?", ids).
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
