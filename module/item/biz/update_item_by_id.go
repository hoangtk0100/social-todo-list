package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error
}

type updateItemBiz struct {
	store     UpdateItemStorage
	requester common.Requester
}

func NewUpdateItemBiz(store UpdateItemStorage, requester common.Requester) *updateItemBiz {
	return &updateItemBiz{
		store:     store,
		requester: requester,
	}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return core.ErrNotFound.
				WithDebug(err.Error())
		}

		return core.ErrInternalServerError.
			WithError(model.ErrCannotGetItem.Error()).
			WithDebug(err.Error())
	}

	if data.Status == "Deleted" {
		return core.ErrBadRequest.
			WithError(model.ErrItemDeleted.Error())
	}

	if !common.IsOwner(biz.requester, data.UserId) && !common.IsAdmin(biz.requester) {
		return core.ErrForbidden.
			WithError(model.ErrRequesterIsNotOwner.Error()).
			WithDebug(err.Error())
	}

	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return core.ErrInternalServerError.
			WithError(model.ErrCannotUpdateItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
