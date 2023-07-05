package business

import (
	"context"

	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

func (biz *business) Login(ctx context.Context, data *entity.UserLogin) (*common.Token, error) {
	user, err := biz.repo.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(entity.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	err = util.CheckPassword(user.Password, common.HashPasswordFormat, user.Salt, data.Password)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(entity.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	user.Mask()
	accessToken, payload, err := biz.tokenMaker.CreateToken(token.AccessToken, user.FakeID.String())
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(entity.ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	tokenResult := &common.Token{
		AccessToken: accessToken,
		ExpiredAt:   payload.ExpiredAt,
	}

	return tokenResult, nil
}
