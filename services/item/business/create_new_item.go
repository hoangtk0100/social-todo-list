package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (biz *business) CreateItem(ctx context.Context, data *entity.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := biz.repo.CreateItem(ctx, data); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateItem.Error()).
			WithDebug(err.Error())
	}

	return nil
}
