package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.delete")
	defer span.End()

	deletedStatus := "Deleted"

	if err := repo.db.Table(entity.TodoItem{}.TableName()).
		Where(cond).
		Updates(map[string]interface{}{
			"status": deletedStatus,
		}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
