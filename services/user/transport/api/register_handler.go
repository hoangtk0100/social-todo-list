package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

func (service *service) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data entity.UserCreate

		if err := ctx.ShouldBind(&data); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		if err := service.business.Register(ctx, &data); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(data.ID))
	}
}
