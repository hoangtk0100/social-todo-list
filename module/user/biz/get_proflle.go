package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

type GetProfileStorage interface {
	GetUserByID(ctx context.Context, id int) (*model.User, error)
}

type getProfileBiz struct {
	store GetProfileStorage
}

func NewGetProfileBiz(store GetProfileStorage) *getProfileBiz {
	return &getProfileBiz{
		store: store,
	}
}

func (biz *getProfileBiz) GetProfile(ctx context.Context) (*model.User, error) {
	requester := core.GetRequester(ctx)
	requesterID := common.GetRequesterID(requester)

	user, err := biz.store.GetUserByID(ctx, requesterID)
	if err != nil {
		return nil, core.ErrUnauthorized.
			WithError(model.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return user, nil
}
