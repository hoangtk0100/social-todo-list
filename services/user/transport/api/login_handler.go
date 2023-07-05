package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

func (service *service) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData entity.UserLogin

		if err := ctx.ShouldBind(&loginData); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		tokenPayload, err := service.business.Login(ctx, &loginData)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(tokenPayload))
	}
}
