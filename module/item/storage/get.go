package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (store *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	_, span := trace.StartSpan(ctx, "item.storage.get")
	defer span.End()

	var data model.TodoItem

	if err := store.db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
