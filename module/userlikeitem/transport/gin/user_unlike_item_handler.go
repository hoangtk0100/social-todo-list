package ginuserlikeitem

import (
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
	"github.com/hoangtk0100/social-todo-list/pubsub"
	"gorm.io/gorm"
)

func UnlikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := ctx.MustGet(common.CurrentUser).(*model.User)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

		store := storage.NewSQLStore(db)
		business := biz.NewUserUnlikeItemBiz(store, ps)

		if err := business.UnlikeItem(ctx.Request.Context(), requester.GetUserId(), int(id.GetLocalID())); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
