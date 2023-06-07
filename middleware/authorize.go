package middleware

import (
	"context"
	"errors"
	"strings"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/plugin/tokenprovider"
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

func RequireAuth(store AuthenStore, serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token, err := extractTokenFromHeader(ctx.GetHeader(authorizationHeaderKey))
		if err != nil {
			panic(err)
		}

		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.TokenProvider)
		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		user, err := store.FindUser(ctx.Request.Context(), map[string]interface{}{"id": payload.UserId()})
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
