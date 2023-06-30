package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	_, span := trace.StartSpan(ctx, "item.storage.insert")
	defer span.End()

	if err := store.db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
