package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"go.opencensus.io/trace"
)

func (store *sqlStore) Delete(ctx context.Context, userId, itemId int) error {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.delete")
	defer span.End()

	if err := store.db.Table(model.Like{}.TableName()).
		Where("item_id = ? and user_id = ?", itemId, userId).
		Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
