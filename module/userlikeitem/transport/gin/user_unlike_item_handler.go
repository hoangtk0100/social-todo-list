package ginuserlikeitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func UnlikeItem(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := ctx.MustGet(common.CurrentUser).(*model.User)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		ps := ac.MustGet(common.PluginPubSub).(core.PubSubComponent)

		store := storage.NewSQLStore(db)
		business := biz.NewUserUnlikeItemBiz(store, ps)

		if err := business.UnlikeItem(ctx.Request.Context(), requester.GetUserId(), int(id.GetLocalID())); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
