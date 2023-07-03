package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type UpdateItemRepository interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *entity.TodoItemUpdate) error
}

type updateItemBusiness struct {
	repo      UpdateItemRepository
	requester core.Requester
}

func NewUpdateItemBusiness(repo UpdateItemRepository, requester core.Requester) *updateItemBusiness {
	return &updateItemBusiness{
		repo:      repo,
		requester: requester,
	}
}

func (biz *updateItemBusiness) UpdateItemByID(ctx context.Context, id int, dataUpdate *entity.TodoItemUpdate) error {
	data, err := biz.repo.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return core.ErrNotFound.
				WithDebug(err.Error())
		}

		return core.ErrInternalServerError.
			WithError(entity.ErrCannotGetItem.Error()).
			WithDebug(err.Error())
	}

	if data.Status == "Deleted" {
		return core.ErrBadRequest.
			WithError(entity.ErrItemDeleted.Error())
	}

	if data.UserID != common.GetRequesterID(biz.requester) {
		return core.ErrForbidden.
			WithError(entity.ErrRequesterIsNotOwner.Error()).
			WithDebug(err.Error())
	}

	if err := biz.repo.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotUpdateItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
