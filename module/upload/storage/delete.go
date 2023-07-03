package storage

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) DeleteImages(ctx context.Context, ids []int) error {
	_, span := trace.StartSpan(ctx, "upload.storage.delete")
	defer span.End()

	if err := store.db.Table(core.Image{}.TableName()).
		Where("id in ?", ids).
		Delete(nil).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
