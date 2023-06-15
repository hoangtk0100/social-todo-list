package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (store *sqlStore) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error) {
	_, span := trace.StartSpan(ctx, "user.storage.get")
	defer span.End()

	db := store.db.Table(model.User{}.TableName())
	for _, value := range moreInfo {
		db = db.Preload(value)
	}

	var user model.User
	if err := db.Where(conds).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &user, nil
}
