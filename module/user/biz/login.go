package biz

import (
	"context"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/tokenprovider"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type loginBiz struct {
	store         LoginStorage
	tokenProvider tokenprovider.TokenProvider
	hasher        Hasher
	expiry        int
}

func NewLoginBiz(store LoginStorage, tokenProvider tokenprovider.TokenProvider, hasher Hasher, expiry int) *loginBiz {
	return &loginBiz{
		store:         store,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

func (biz *loginBiz) Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error) {
	user, err := biz.store.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	hashedPassword := biz.hasher.Hash(data.Password + user.Salt)
	if user.Password != hashedPassword {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	payload := &common.TokenPayload{
		UId:   user.Id,
		URole: user.Role.String(),
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil
}
