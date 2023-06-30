package biz

import (
	"context"
	"strconv"

	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type loginBiz struct {
	store      LoginStorage
	tokenMaker core.TokenMakerComponent
}

func NewLoginBiz(store LoginStorage, tokenMaker core.TokenMakerComponent) *loginBiz {
	return &loginBiz{
		store:      store,
		tokenMaker: tokenMaker,
	}
}

func (biz *loginBiz) Login(ctx context.Context, data *model.UserLogin) (*common.Token, error) {
	user, err := biz.store.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(model.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	err = util.CheckPassword(user.Password, common.HashPasswordFormat, user.Salt, data.Password)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(model.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	accessToken, payload, err := biz.tokenMaker.CreateToken(token.AccessToken, strconv.Itoa(user.Id))
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(model.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	tokenResult := &common.Token{
		AccessToken: accessToken,
		ExpiredAt:   payload.ExpiredAt,
	}

	return tokenResult, nil
}
