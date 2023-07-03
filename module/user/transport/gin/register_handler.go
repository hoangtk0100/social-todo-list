package ginuser

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/biz"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/user/storage"
)

func Register(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data model.UserCreate

		if err := ctx.ShouldBind(&data); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		store := storage.NewSQLStore(db)
		business := biz.NewRegisterBiz(store)

		if err := business.Register(ctx.Request.Context(), &data); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(data.ID))
	}
}
