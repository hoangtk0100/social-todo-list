package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
)

func (store *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := store.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
