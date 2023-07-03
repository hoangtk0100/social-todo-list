package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type DeleteItemRepository interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBusiness struct {
	repo DeleteItemRepository
}

func NewDeleteItemBusiness(repo DeleteItemRepository) *deleteItemBusiness {
	return &deleteItemBusiness{repo: repo}
}

func (biz *deleteItemBusiness) DeleteItemByID(ctx context.Context, id int) error {
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

	if err := biz.repo.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotDeleteItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
