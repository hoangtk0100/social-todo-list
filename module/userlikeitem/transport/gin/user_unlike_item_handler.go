package ginuserlikeitem

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	usermodel "github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func UnlikeItem(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(model.ErrItemIdInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		requester := ctx.MustGet(common.CurrentUser).(*usermodel.User)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		ps := ac.MustGet(common.PluginPubSub).(core.PubSubComponent)

		store := storage.NewSQLStore(db)
		business := biz.NewUserUnlikeItemBiz(store, ps)

		if err := business.UnlikeItem(ctx.Request.Context(), requester.GetUserId(), int(id.GetLocalID())); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(true))
	}
}
