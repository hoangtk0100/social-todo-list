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

func UpdateItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(model.ErrItemIdInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		var data model.TodoItemUpdate
		if err := ctx.ShouldBind(&data); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		store := storage.NewSQLStore(db)
		business := biz.NewUpdateItemBiz(store, requester)

		if err := business.UpdateItemById(ctx.Request.Context(), int(id.GetLocalID()), &data); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(true))
	}
}
