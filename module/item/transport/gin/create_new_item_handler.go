package ginitem

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/biz"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
)

func CreateItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var itemData model.TodoItemCreation

		if err := ctx.ShouldBind(&itemData); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		requester := core.GetRequester(ctx)
		itemData.UserID = common.GetRequesterID(requester)

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(ctx.Request.Context(), &itemData); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(itemData.ID))
	}
}
