package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

type UserRepository interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, data *entity.UserCreate) error
}

type business struct {
	repo       UserRepository
	tokenMaker core.TokenMakerComponent
}

func NewBusiness(repo UserRepository, tokenMaker core.TokenMakerComponent) *business {
	return &business{
		repo:       repo,
		tokenMaker: tokenMaker,
	}
}
