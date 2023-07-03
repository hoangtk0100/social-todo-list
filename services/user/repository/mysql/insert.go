package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/user/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) CreateUser(ctx context.Context, data *entity.UserCreate) error {
	_, span := trace.StartSpan(ctx, "user.repository.mysql.insert")
	defer span.End()

	db := repo.db.Begin()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	return nil
}
