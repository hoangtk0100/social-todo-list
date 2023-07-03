package api

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/business"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
)

func CreateItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var itemData entity.TodoItemCreation

		if err := ctx.ShouldBind(&itemData); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		requester := core.GetRequester(ctx)
		itemData.UserID = common.GetRequesterID(requester)

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		repo := mysql.NewMySQLRepository(db)
		business := business.NewCreateItemBusiness(repo)

		if err := business.CreateNewItem(ctx.Request.Context(), &itemData); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(itemData.ID))
	}
}
