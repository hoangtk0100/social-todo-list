package rpcuserlikeitem

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func GetItemLikes(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type requestData struct {
			IDs []int `json:"ids"`
		}

		var data requestData
		if err := ctx.ShouldBind(&data); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		store := storage.NewSQLStore(db)

		mapResult, err := store.GetItemLikes(ctx.Request.Context(), data.IDs)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(mapResult))
	}
}
