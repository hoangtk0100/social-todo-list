package mysql

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (repo *mysqlRepo) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*entity.User, error) {
	_, span := trace.StartSpan(ctx, "user.repository.mysql.find")
	defer span.End()

	db := repo.db.Table(entity.User{}.TableName())
	for _, value := range moreInfo {
		db = db.Preload(value)
	}

	var user entity.User
	if err := db.Where(conds).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}

func (repo *mysqlRepo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	_, span := trace.StartSpan(ctx, "user.repository.mysql.get_details")
	defer span.End()

	db := repo.db.Table(entity.User{}.TableName())

	var user entity.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}
