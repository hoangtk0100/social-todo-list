package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

func (store *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := store.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return err
	}

	return nil
}
