package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (biz *business) ListItem(ctx context.Context,
	filter *entity.Filter,
	paging *core.Paging,
) ([]entity.TodoItem, error) {
	data, err := biz.repo.ListItem(ctx, filter, paging, "Owner")
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotGetItems.Error()).
			WithDebug(err.Error())
	}

	if len(data) == 0 {
		return data, nil
	}

	ids := make([]int, len(data))
	for index := range ids {
		ids[index] = data[index].ID
	}

	itemLikesMap, err := biz.likeRepo.GetItemLikes(ctx, ids)
	if err != nil {
		return data, nil
	}

	for index := range data {
		data[index].LikedCount = itemLikesMap[data[index].ID]
	}

	return data, nil
}
