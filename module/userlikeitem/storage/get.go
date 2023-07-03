package storage

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (store *sqlStore) Find(ctx context.Context, userID, itemID int) (*model.Like, error) {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.get")
	defer span.End()

	var data model.Like

	if err := store.db.
		Where("item_id = ? and user_id = ?", itemID, userID).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &data, nil
}
