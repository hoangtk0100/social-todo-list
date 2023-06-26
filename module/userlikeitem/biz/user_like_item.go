package biz

import (
	"context"
	"log"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
)

type UserLikeItemStorage interface {
	Create(ctx context.Context, data *model.Like) error
}

type userLikeItemBiz struct {
	store UserLikeItemStorage
	ps    core.PubSubComponent
}

func NewUserLikeItemBiz(store UserLikeItemStorage, ps core.PubSubComponent) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, ps: ps}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	if err := biz.ps.Publish(ctx, common.TopicUserLikedItem, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	return nil
}
