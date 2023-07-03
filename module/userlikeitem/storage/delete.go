package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (store *sqlStore) Delete(ctx context.Context, userID, itemID int) error {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.delete")
	defer span.End()

	if err := store.db.Table(model.Like{}.TableName()).
		Where("item_id = ? and user_id = ?", itemID, userID).
		Delete(nil).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
