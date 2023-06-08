package biz

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
)

type UserLikeItemStorage interface {
	Create(ctx context.Context, data *model.Like) error
}

type userLikeItemBiz struct {
	store UserLikeItemStorage
}

func NewUserLikeItemBiz(store UserLikeItemStorage) *userLikeItemBiz {
	return &userLikeItemBiz{store: store}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	return nil
}
