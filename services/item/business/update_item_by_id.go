package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (biz *business) UpdateItemByID(ctx context.Context, id int, dataUpdate *entity.TodoItemUpdate) error {
	data, err := biz.repo.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return core.ErrNotFound.
				WithError(entity.ErrCannotGetItem.Error()).
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

	requester := core.GetRequester(ctx)
	requesterID := common.GetRequesterID(requester)
	if data.UserID != requesterID {
		return core.ErrForbidden.
			WithError(entity.ErrRequesterIsNotOwner.Error())
	}

	if err := biz.repo.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotUpdateItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
