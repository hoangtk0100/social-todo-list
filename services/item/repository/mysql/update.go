package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (repo *mysqlRepo) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *entity.TodoItemUpdate) error {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.update")
	defer span.End()

	if err := repo.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *mysqlRepo) IncreaseLikedCount(ctx context.Context, id int) error {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.increase_liked_count")
	defer span.End()

	if err := repo.db.Table(entity.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + ?", 1)).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *mysqlRepo) DecreaseLikedCount(ctx context.Context, id int) error {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.decrease_liked_count")
	defer span.End()

	if err := repo.db.Table(entity.TodoItem{}.TableName()).
		Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count - ?", 1)).
		Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
