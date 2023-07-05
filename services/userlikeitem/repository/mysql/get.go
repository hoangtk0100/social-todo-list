package mysql

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (repo *mysqlRepo) Find(ctx context.Context, userID, itemID int) (*entity.Like, error) {
	_, span := trace.StartSpan(ctx, "userlikeitem.repository.mysql.get")
	defer span.End()

	var data entity.Like

	if err := repo.db.
		Where("item_id = ? and user_id = ?", itemID, userID).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &data, nil
}
