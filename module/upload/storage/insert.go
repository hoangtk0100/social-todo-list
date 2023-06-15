package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"go.opencensus.io/trace"
)

func (store *sqlStore) CreateImage(ctx context.Context, data *common.Image) error {
	_, span := trace.StartSpan(ctx, "upload.storage.insert")
	defer span.End()

	if err := store.db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
