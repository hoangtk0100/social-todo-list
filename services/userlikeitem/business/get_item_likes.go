package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

func (biz *business) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result, err := biz.repo.GetItemLikes(ctx, ids)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotGetItemLikes.Error()).
			WithDebug(err.Error())
	}

	return result, nil
}
