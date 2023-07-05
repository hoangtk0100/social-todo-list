package api

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

type UserLikeItemBusiness interface {
	LikeItem(ctx context.Context, data *entity.Like) error
	UnlikeItem(ctx context.Context, userID, itemID int) error
	ListUsersLikedItem(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error)
}

type service struct {
	ac       appctx.AppContext
	business UserLikeItemBusiness
}

func NewService(ac appctx.AppContext, business UserLikeItemBusiness) *service {
	return &service{
		ac:       ac,
		business: business,
	}
}
