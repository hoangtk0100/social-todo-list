package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

func (store *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := store.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
