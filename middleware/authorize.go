package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
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

type AuthenStore interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
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

func RequireAuth(store AuthenStore, ac appctx.AppContext) gin.HandlerFunc {
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

		userId, _ := strconv.Atoi(payload.UID)
		user, err := store.FindUser(ctx.Request.Context(), map[string]interface{}{"id": userId})
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

		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
