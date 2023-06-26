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

func Login(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData model.UserLogin

		if err := ctx.ShouldBind(&loginData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		tokenMaker := ac.MustGet(common.PluginJWT).(core.TokenMakerComponent)

		store := storage.NewSQLStore(db)
		business := biz.NewLoginBiz(store, tokenMaker)
		tokenPayload, err := business.Login(ctx.Request.Context(), &loginData)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(tokenPayload))
	}
}
