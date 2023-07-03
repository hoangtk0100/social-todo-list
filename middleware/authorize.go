package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
)

var (
	ErrAuthHeaderEmpty           = errors.New("authorization header is not provided")
	ErrAuthHeaderInvalidFormat   = errors.New("invalid authorization header format")
	ErrAuthHeaderUnsupportedType = errors.New("unsupported authorization type")
	ErrUserDeletedOrBanned       = errors.New("user has been deleted or banned")
)

type AuthenRepository interface {
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

func extractTokenFromHeader(input string) (string, error) {
	if len(input) == 0 {
		return "", ErrAuthHeaderEmpty
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "", ErrAuthHeaderInvalidFormat
	}

	authorizationType := strings.ToLower(parts[0])
	if authorizationType != authorizationTypeBearer {
		return "", ErrAuthHeaderUnsupportedType
	}

	return parts[1], nil
}

func RequireAuth(repo AuthenRepository, ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extractTokenFromHeader(ctx.GetHeader(authorizationHeaderKey))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrUnauthorized.WithError(err.Error()))
			ctx.Abort()
			return
		}

		tokenMaker := ac.MustGet(common.PluginJWT).(core.TokenMakerComponent)
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			core.ErrorResponse(ctx, core.ErrUnauthorized.WithError(err.Error()).WithDebug(err.Error()))
			ctx.Abort()
			return
		}

		uid, err := util.UIDFromString(payload.UID)
		if err != nil {
			core.ErrorResponse(ctx, core.ErrUnauthorized.WithError(err.Error()).WithDebug(err.Error()))
			ctx.Abort()
			return
		}

		userID := int(uid.GetLocalID())
		user, err := repo.GetUserByID(ctx.Request.Context(), userID)
		if err != nil {
			core.ErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if user.Status == 0 {
			core.ErrorResponse(ctx, core.ErrForbidden.WithError(ErrUserDeletedOrBanned.Error()))
			ctx.Abort()
			return
		}

		ctx.Set(core.KeyRequester, core.NewRequester(payload.ID.String(), payload.UID))
		ctx.Next()
	}
}
