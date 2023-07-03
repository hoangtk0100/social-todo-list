package storage

import (
	"context"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const (
	timeLayout   = "2006-01-02T15:04:05.999999"
	timeDBLayout = "2006-01-02 15:04:05.999999"
)

func (store *sqlStore) ListUsers(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error) {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.list_users")
	defer span.End()

	var result []model.Like
	db := store.db.Table(model.Like{}.TableName()).Where("item_id = ?", itemID)

	if err := db.Select("user_id").Count(&paging.Total).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	if cursor := strings.TrimSpace(paging.FakeCursor); cursor != "" {
		createdTime, err := time.Parse(timeLayout, string(base58.Decode(cursor)))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		db = db.Where("created_at < ?", createdTime.Format(timeDBLayout))
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").
		Order("created_at desc").
		Limit(paging.Limit).
		Preload("User").
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	size := len(result)
	users := make([]core.SimpleUser, size)
	for index := range users {
		users[index] = *result[index].User
		users[index].UpdatedAt = nil
		users[index].CreatedAt = result[index].CreatedAt
	}

	if size > 0 {
		paging.NextCursor = base58.Encode([]byte(users[size-1].CreatedAt.Format(timeLayout)))
	}

	return users, nil
}

func (store *sqlStore) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	_, span := trace.StartSpan(ctx, "userlikeitem.storage.get_item_likes")
	defer span.End()

	result := make(map[int]int)

	type sqlData struct {
		ItemID int `gorm:"column:item_id"`
		Count  int `gorm:"column:count"`
	}

	var likes []sqlData
	if err := store.db.Table(model.Like{}.TableName()).
		Select("item_id, Count(item_id) as `count`").
		Where("item_id in (?)", ids).
		Group("item_id").
		Find(&likes).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	for _, item := range likes {
		result[item.ItemID] = item.Count
	}

	return result, nil
}
