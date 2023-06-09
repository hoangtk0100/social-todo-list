package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"gorm.io/gorm"
)

func (store *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := store.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (store *sqlStore) IncreaseLikedCount(ctx context.Context, id int) error {
	if err := store.db.Table(model.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + ?", 1)).
		Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (store *sqlStore) DecreaseLikedCount(ctx context.Context, id int) error {
	if err := store.db.Table(model.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count - ?", 1)).
		Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
