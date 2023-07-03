package storage

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) CreateImage(ctx context.Context, data *core.Image) error {
	_, span := trace.StartSpan(ctx, "upload.storage.insert")
	defer span.End()

	if err := store.db.Table(data.TableName()).Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
