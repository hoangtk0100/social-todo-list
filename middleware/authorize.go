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

func ErrAuthHeaderEmpty(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"authorization header is not provided",
		"ErrAuthHeaderEmpty",
	)
}

func ErrAuthHeaderInvalidFormat(err error) *common.AppError {
	return common.NewUnauthorized(
		err,
		"invalid authorization header format",
		"ErrAuthHeaderInvalidFormat",
	)
}

func ErrAuthHeaderUnsupportedType(err error) *common.AppError {
	return common.NewUnauthorized(
		err,
		"unsupported authorization type",
		"ErrAuthHeaderUnsupportedType",
	)
}

type AuthenStore interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
}

func extractTokenFromHeader(input string) (string, error) {
	if len(input) == 0 {
		return "", ErrAuthHeaderEmpty(nil)
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "", ErrAuthHeaderInvalidFormat(nil)
	}

	authorizationType := strings.ToLower(parts[0])
	if authorizationType != authorizationTypeBearer {
		return "", ErrAuthHeaderUnsupportedType(nil)
	}

	return parts[1], nil
}

func RequireAuth(store AuthenStore, ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extractTokenFromHeader(ctx.GetHeader(authorizationHeaderKey))
		if err != nil {
			panic(err)
		}

		tokenMaker := ac.MustGet(common.PluginJWT).(core.TokenMakerComponent)
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			panic(err)
		}

		userId, _ := strconv.Atoi(payload.UID)
		user, err := store.FindUser(ctx.Request.Context(), map[string]interface{}{"id": userId})
		if err != nil {
			panic(err)
		}

		if user.Status == 0 {
			panic(common.ErrNoPermission(errors.New("user has been deleted or banned")))
		}

		ctx.Set(common.CurrentUser, user)
		ctx.Next()
	}
}
