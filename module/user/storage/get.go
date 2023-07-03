package storage

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (store *sqlStore) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error) {
	_, span := trace.StartSpan(ctx, "user.storage.find")
	defer span.End()

	db := store.db.Table(model.User{}.TableName())
	for _, value := range moreInfo {
		db = db.Preload(value)
	}

	var user model.User
	if err := db.Where(conds).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}

func (store *sqlStore) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	_, span := trace.StartSpan(ctx, "user.storage.get_details")
	defer span.End()

	db := store.db.Table(model.User{}.TableName())

	var user model.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}
