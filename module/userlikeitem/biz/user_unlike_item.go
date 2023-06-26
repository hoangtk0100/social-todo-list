package biz

import (
	"context"
	"log"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
)

type UserUnlikeItemStorage interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type userUnlikeItemBiz struct {
	store UserUnlikeItemStorage
	ps    core.PubSubComponent
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStorage, ps core.PubSubComponent) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, ps: ps}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)
	if err == common.RecordNotFound {
		return model.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	msg := pubsub.NewMessage(&model.Like{
		UserId: userId,
		ItemId: itemId,
	})

	if err := biz.ps.Publish(ctx, common.TopicUserUnlikedItem, msg); err != nil {
		log.Println(err)
	}

	return nil
}
