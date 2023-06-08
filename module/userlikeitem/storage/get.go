package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"gorm.io/gorm"
)

func (store *sqlStore) Find(ctx context.Context, userId, itemId int) (*model.Like, error) {
	var data model.Like

	if err := store.db.
		Where("item_id = ? and user_id = ?", itemId, userId).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
