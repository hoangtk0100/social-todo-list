package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/rs/zerolog/log"
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
		return core.ErrInternalServerError.
			WithError(model.ErrCannotLikeItem.Error()).
			WithDebug(err.Error())
	}

	msg := pubsub.NewMessage(data)
	if err := biz.ps.Publish(ctx, common.TopicUserLikedItem, msg); err != nil {
		log.Error().Err(err).Msgf(
			"Topic(%s) - Msg(%s) - Error(%s)",
			common.TopicUserLikedItem,
			msg.String(),
			err.Error(),
		)
	}

	return nil
}
