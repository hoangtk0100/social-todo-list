package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, data *entity.TodoItemCreation) error
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *entity.TodoItemUpdate) error
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
	ListItem(
		ctx context.Context,
		filter *entity.Filter,
		paging *core.Paging,
		moreKeys ...string,
	) ([]entity.TodoItem, error)
}

type ItemLikeRepository interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type business struct {
	repo     ItemRepository
	likeRepo ItemLikeRepository
}

func NewBusiness(repo ItemRepository, likeRepo ItemLikeRepository) *business {
	return &business{
		repo:     repo,
		likeRepo: likeRepo,
	}
}
