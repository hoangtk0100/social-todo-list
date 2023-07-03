package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
}

type getItemBiz struct {
	store GetItemStorage
}

func NewGetItemBiz(store GetItemStorage) *getItemBiz {
	return &getItemBiz{store: store}
}

func (biz *getItemBiz) GetItemByID(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.
				WithDebug(err.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(model.ErrCannotGetItem.Error()).
			WithDebug(err.Error())
	}

	return data, nil
}
