package business

import (
	"context"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/rs/zerolog/log"
)

func (biz *business) LikeItem(ctx context.Context, data *entity.Like) error {
	if err := biz.repo.Create(ctx, data); err != nil {
		return core.ErrConflict.
			WithError(entity.ErrLikedItem.Error()).
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
