package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

func (store *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem

	if err := store.db.Where(cond).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
