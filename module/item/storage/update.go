package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (store *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	_, span := trace.StartSpan(ctx, "item.storage.update")
	defer span.End()

	if err := store.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *sqlStore) IncreaseLikedCount(ctx context.Context, id int) error {
	_, span := trace.StartSpan(ctx, "item.storage.increase_liked_count")
	defer span.End()

	if err := store.db.Table(model.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + ?", 1)).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *sqlStore) DecreaseLikedCount(ctx context.Context, id int) error {
	_, span := trace.StartSpan(ctx, "item.storage.decrease_liked_count")
	defer span.End()

	if err := store.db.Table(model.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count - ?", 1)).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
