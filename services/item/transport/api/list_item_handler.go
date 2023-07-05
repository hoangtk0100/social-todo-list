package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (service *service) ListItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queryString struct {
			core.Paging
			entity.Filter
		}

		if err := ctx.ShouldBind(&queryString); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		queryString.Paging.Process()

		result, err := service.business.ListItem(ctx, &queryString.Filter, &queryString.Paging)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		for index := range result {
			result[index].Mask()
		}

		core.SuccessResponse(ctx, core.NewResponse(result, queryString.Paging, queryString.Filter))
	}
}
