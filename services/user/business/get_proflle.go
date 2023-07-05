package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

func (biz *business) GetProfile(ctx context.Context) (*entity.User, error) {
	requester := core.GetRequester(ctx)
	requesterID := common.GetRequesterID(requester)

	user, err := biz.repo.GetUserByID(ctx, requesterID)
	if err != nil {
		return nil, core.ErrUnauthorized.
			WithError(entity.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return user, nil
}
