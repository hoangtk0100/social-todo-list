package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

type ListItemRepository interface {
	ListItem(
		ctx context.Context,
		filter *entity.Filter,
		paging *core.Paging,
		moreKeys ...string,
	) ([]entity.TodoItem, error)
}

type ItemLikeRepository interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listItemBusiness struct {
	repo      ListItemRepository
	likeRepo  ItemLikeRepository
	requester core.Requester
}

func NewListItemBusiness(repo ListItemRepository, likeRepo ItemLikeRepository, requester core.Requester) *listItemBusiness {
	return &listItemBusiness{
		repo:      repo,
		likeRepo:  likeRepo,
		requester: requester,
	}
}

func (biz *listItemBusiness) ListItem(ctx context.Context,
	filter *entity.Filter,
	paging *core.Paging,
) ([]entity.TodoItem, error) {
	newCtx := core.ContextWithRequester(ctx, biz.requester)

	data, err := biz.repo.ListItem(newCtx, filter, paging, "Owner")
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

	itemLikesMap, err := biz.likeRepo.GetItemLikes(newCtx, ids)
	if err != nil {
		return data, nil
	}

	for index := range data {
		data[index].LikedCount = itemLikesMap[data[index].ID]
	}

	return data, nil
}
