package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type GetItemRepository interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
}

type getItemBusiness struct {
	repo GetItemRepository
}

func NewGetItemBusiness(repo GetItemRepository) *getItemBusiness {
	return &getItemBusiness{repo: repo}
}

func (biz *getItemBusiness) GetItemByID(ctx context.Context, id int) (*entity.TodoItem, error) {
	data, err := biz.repo.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.
				WithDebug(err.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotGetItem.Error()).
			WithDebug(err.Error())
	}

	return data, nil
}
