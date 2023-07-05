package api

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type Business interface {
	CreateItem(ctx context.Context, data *entity.TodoItemCreation) error
	DeleteItemByID(ctx context.Context, id int) error
	GetItemByID(ctx context.Context, id int) (*entity.TodoItem, error)
	UpdateItemByID(ctx context.Context, id int, dataUpdate *entity.TodoItemUpdate) error
	ListItem(ctx context.Context,
		filter *entity.Filter,
		paging *core.Paging,
	) ([]entity.TodoItem, error)
}

type service struct {
	ac       appctx.AppContext
	business Business
}

func NewService(ac appctx.AppContext, business Business) *service {
	return &service{
		ac:       ac,
		business: business,
	}
}
