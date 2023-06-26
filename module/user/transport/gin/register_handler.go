package ginuser

import (
	"net/http"

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
			panic(common.ErrInvalidRequest(err))
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		business := biz.NewRegisterBiz(store, md5)

		if err := business.Register(ctx.Request.Context(), &data); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
