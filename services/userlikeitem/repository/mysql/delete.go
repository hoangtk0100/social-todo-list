package mysql

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) Delete(ctx context.Context, userID, itemID int) error {
	_, span := trace.StartSpan(ctx, "userlikeitem.repository.mysql.delete")
	defer span.End()

	if err := repo.db.Table(entity.Like{}.TableName()).
		Where("item_id = ? and user_id = ?", itemID, userID).
		Delete(nil).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
