package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

type ListUsersLikedItemRepository interface {
	ListUsers(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error)
}

type listUsersLikedItemBusiness struct {
	repo ListUsersLikedItemRepository
}

func NewListUsersLikedItemBusiness(repo ListUsersLikedItemRepository) *listUsersLikedItemBusiness {
	return &listUsersLikedItemBusiness{repo: repo}
}

func (biz *listUsersLikedItemBusiness) ListUsersLikedItem(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error) {
	result, err := biz.repo.ListUsers(ctx, itemID, paging)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotListLikedUsers.Error()).
			WithDebug(err.Error())
	}

	return result, nil
}
