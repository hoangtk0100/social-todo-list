package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type CreateItemRepository interface {
	CreateItem(ctx context.Context, data *entity.TodoItemCreation) error
}

type createItemBusiness struct {
	repo CreateItemRepository
}

func NewCreateItemBusiness(repo CreateItemRepository) *createItemBusiness {
	return &createItemBusiness{repo: repo}
}

func (biz *createItemBusiness) CreateNewItem(ctx context.Context, data *entity.TodoItemCreation) error {
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
