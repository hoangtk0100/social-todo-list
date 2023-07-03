package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) CreateItem(ctx context.Context, data *entity.TodoItemCreation) error {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.insert")
	defer span.End()

	if err := repo.db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
