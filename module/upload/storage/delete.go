package storage

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
)

func (store *sqlStore) DeleteImages(ctx context.Context, ids []int) error {
	if err := store.db.Table(common.Image{}.TableName()).
		Where("id in ?", ids).
		Delete(nil).Error; err != nil {
		return err
	}

	return nil
}
