package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"go.opencensus.io/trace"
)

func (store *sqlStore) CreateUser(ctx context.Context, data *model.UserCreate) error {
	_, span := trace.StartSpan(ctx, "user.storage.insert")
	defer span.End()

	db := store.db.Begin()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
