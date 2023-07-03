package business

import (
	"context"

	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/rs/zerolog/log"
)

type UserUnlikeItemRepository interface {
	Find(ctx context.Context, userID, itemID int) (*entity.Like, error)
	Delete(ctx context.Context, userID, itemID int) error
}

type userUnlikeItemBusiness struct {
	repo UserUnlikeItemRepository
	ps   core.PubSubComponent
}

func NewUserUnlikeItemBusiness(repo UserUnlikeItemRepository, ps core.PubSubComponent) *userUnlikeItemBusiness {
	return &userUnlikeItemBusiness{repo: repo, ps: ps}
}

func (biz *userUnlikeItemBusiness) UnlikeItem(ctx context.Context, userID, itemID int) error {
	_, err := biz.repo.Find(ctx, userID, itemID)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return core.ErrNotFound.
				WithError(entity.ErrDidNotLikeItem.Error()).
				WithDebug(err.Error())
		}

		return core.ErrInternalServerError.
			WithError(entity.ErrCannotUnlikeItem.Error()).
			WithDebug(err.Error())
	}

	if err := biz.repo.Delete(ctx, userID, itemID); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotUnlikeItem.Error()).
			WithDebug(err.Error())
	}

	msg := pubsub.NewMessage(&entity.Like{
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
