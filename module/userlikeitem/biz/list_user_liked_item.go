package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
)

type ListUsersLikedItemStorage interface {
	ListUsers(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error)
}

type listUsersLikedItemBiz struct {
	store ListUsersLikedItemStorage
}

func NewListUsersLikedItemBiz(store ListUsersLikedItemStorage) *listUsersLikedItemBiz {
	return &listUsersLikedItemBiz{store: store}
}

func (biz *listUsersLikedItemBiz) ListUsersLikedItem(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error) {
	result, err := biz.store.ListUsers(ctx, itemId, paging)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(model.ErrCannotListLikedUsers.Error()).
			WithDebug(err.Error())
	}

	return result, nil
}
