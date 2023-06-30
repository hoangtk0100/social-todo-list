package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) DeleteImages(ctx context.Context, ids []int) error {
	_, span := trace.StartSpan(ctx, "upload.storage.delete")
	defer span.End()

	if err := store.db.Table(common.Image{}.TableName()).
		Where("id in ?", ids).
		Delete(nil).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
