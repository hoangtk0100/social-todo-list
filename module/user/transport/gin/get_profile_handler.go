package ginuser

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/biz"
	"github.com/hoangtk0100/social-todo-list/module/user/storage"
)

func Profile(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		store := storage.NewSQLStore(db)
		business := biz.NewGetProfileBiz(store)

		user, err := business.GetProfile(ctx)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		data := core.NewSimpleUser(user.ID, user.FirstName, user.LastName, user.Avatar)
		data.Mask(common.MaskTypeUser)
		core.SuccessResponse(ctx, core.NewDataResponse(data))
	}
}
