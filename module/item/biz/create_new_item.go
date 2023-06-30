package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
}

type createItemBiz struct {
	store CreateItemStorage
}

func NewCreateItemBiz(store CreateItemStorage) *createItemBiz {
	return &createItemBiz{store: store}
}

func (biz *createItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := biz.store.CreateItem(ctx, data); err != nil {
		return core.ErrInternalServerError.
			WithError(model.ErrCannotCreateItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
