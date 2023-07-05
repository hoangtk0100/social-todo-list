package mysql

import (
	"strings"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"
)

func (repo *mysqlRepo) ListItem(
	ctx context.Context,
	filter *entity.Filter,
	paging *core.Paging,
	moreKeys ...string,
) ([]entity.TodoItem, error) {
	_, span := trace.StartSpan(ctx, "item.repository.mysql.list")
	defer span.End()

	var result []entity.TodoItem

	db := repo.db.
		Table(entity.TodoItem{}.TableName()).
		Where("status <> ?", "Deleted")

	// Get items of requester only
	// requester := ctx.Value(core.KeyRequester).(core.Requester)
	// db = db.Where("user_id = ?", requester.GetUID())

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	for _, value := range moreKeys {
		db = db.Preload(value)
	}

	if cursor := strings.TrimSpace(paging.FakeCursor); cursor != "" {
		id, err := util.UIDFromString(cursor)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		db = db.Where("id < ?", id.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Select("*").
		Order("id desc").
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	size := len(result)
	if size > 0 {
		result[size-1].Mask()
		paging.NextCursor = result[size-1].FakeID.String()
	}

	return result, nil
}
