package ginitem

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/biz"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
)

func DeleteItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(model.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		store := storage.NewSQLStore(db)
		business := biz.NewDeleteItemBiz(store)

		if err := business.DeleteItemByID(ctx.Request.Context(), int(id.GetLocalID())); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(true))
	}
}
