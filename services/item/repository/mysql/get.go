package mysql

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (repo *mysqlRepo) GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error) {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.get")
	defer span.End()

	var data entity.TodoItem

	if err := repo.db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &data, nil
}
