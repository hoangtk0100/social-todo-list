package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (service *service) CreateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var itemData entity.TodoItemCreation

		if err := ctx.ShouldBind(&itemData); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		requester := core.GetRequester(ctx)
		itemData.UserID = common.GetRequesterID(requester)

		if err := service.business.CreateItem(ctx, &itemData); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(itemData.ID))
	}
}
