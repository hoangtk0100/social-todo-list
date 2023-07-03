package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) Create(ctx context.Context, data *entity.Like) error {
	_, span := trace.StartSpan(ctx, "userlikeitem.repository.mysql.insert")
	defer span.End()

	if err := repo.db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
