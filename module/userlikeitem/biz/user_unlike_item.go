package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/rs/zerolog/log"
)

type UserUnlikeItemStorage interface {
	Find(ctx context.Context, userID, itemID int) (*model.Like, error)
	Delete(ctx context.Context, userID, itemID int) error
}

type userUnlikeItemBiz struct {
	store UserUnlikeItemStorage
	ps    core.PubSubComponent
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStorage, ps core.PubSubComponent) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, ps: ps}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userID, itemID int) error {
	_, err := biz.store.Find(ctx, userID, itemID)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return core.ErrNotFound.
				WithError(model.ErrDidNotLikeItem.Error()).
				WithDebug(err.Error())
		}

		return core.ErrInternalServerError.
			WithError(model.ErrCannotUnlikeItem.Error()).
			WithDebug(err.Error())
	}

	if err := biz.store.Delete(ctx, userID, itemID); err != nil {
		return core.ErrInternalServerError.
			WithError(model.ErrCannotUnlikeItem.Error()).
			WithDebug(err.Error())
	}

	msg := pubsub.NewMessage(&model.Like{
		UserID: userID,
		ItemID: itemID,
	})

	if err := biz.ps.Publish(ctx, common.TopicUserUnlikedItem, msg); err != nil {
		log.Error().Err(err).Msgf(
			"Topic(%s) - Msg(%s) - Error(%s)",
			common.TopicUserUnlikedItem,
			msg.String(),
			err.Error(),
		)
	}

	return nil
}
