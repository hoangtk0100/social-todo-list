package api

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

type Business interface {
	Register(ctx context.Context, data *entity.UserCreate) error
	Login(ctx context.Context, data *entity.UserLogin) (*common.Token, error)
	GetProfile(ctx context.Context) (*entity.User, error)
}

type service struct {
	ac       appctx.AppContext
	business Business
}

func NewService(ac appctx.AppContext, business Business) *service {
	return &service{
		ac:       ac,
		business: business,
	}
}
