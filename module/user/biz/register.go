package biz

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreate) error
}

type registerBiz struct {
	store RegisterStorage
}

func NewRegisterBiz(store RegisterStorage) *registerBiz {
	return &registerBiz{
		store: store,
	}
}

func (biz *registerBiz) Register(ctx context.Context, data *model.UserCreate) error {
	user, _ := biz.store.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return core.ErrBadRequest.
			WithError(model.ErrEmailExisted.Error())
	}

	var err error
	salt := common.GenSalt(50)
	data.Password, err = util.HashPassword(common.HashPasswordFormat, salt, data.Password)
	if err != nil {
		return core.ErrInternalServerError.
			WithError(model.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	data.Salt = salt
	data.Role = model.RoleUser.String()

	if err := biz.store.CreateUser(ctx, data); err != nil {
		return core.ErrInternalServerError.
			WithError(model.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	return nil
}
