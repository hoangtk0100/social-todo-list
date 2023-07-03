package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

type RegisterRepository interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*entity.User, error)
	CreateUser(ctx context.Context, data *entity.UserCreate) error
}

type registerBusiness struct {
	repo RegisterRepository
}

func NewRegisterBusiness(repo RegisterRepository) *registerBusiness {
	return &registerBusiness{
		repo: repo,
	}
}

func (biz *registerBusiness) Register(ctx context.Context, data *entity.UserCreate) error {
	user, _ := biz.repo.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return core.ErrBadRequest.
			WithError(entity.ErrEmailExisted.Error())
	}

	var err error
	salt, err := util.RandomString(50)
	if err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	data.Password, err = util.HashPassword(common.HashPasswordFormat, salt, data.Password)
	if err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	data.Salt = salt
	data.Role = entity.RoleUser.String()

	if err := biz.repo.CreateUser(ctx, data); err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	return nil
}
