package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"go.opencensus.io/trace"
)

func (store *sqlStore) Create(ctx context.Context, data *model.Like) error {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.insert")
	defer span.End()

	if err := store.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
