package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	_, span := trace.StartSpan(ctx, "item.storage.delete")
	defer span.End()

	deletedStatus := "Deleted"

	if err := store.db.Table(model.TodoItem{}.TableName()).
		Where(cond).
		Updates(map[string]interface{}{
			"status": deletedStatus,
		}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
