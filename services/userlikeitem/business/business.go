package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

type UserLikeItemRepository interface {
	Create(ctx context.Context, data *entity.Like) error
	Delete(ctx context.Context, userID, itemID int) error
	Find(ctx context.Context, userID, itemID int) (*entity.Like, error)
	ListUsers(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error)
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type business struct {
	repo UserLikeItemRepository
	ps   core.PubSubComponent
}

func NewBusiness(repo UserLikeItemRepository, ps core.PubSubComponent) *business {
	return &business{
		repo: repo,
		ps:   ps,
	}
}
