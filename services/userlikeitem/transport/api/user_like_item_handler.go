package api

import (
	"time"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/business"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/repository/mysql"
)

func LikeItem(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(entity.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		requester := core.GetRequester(ctx)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		ps := ac.MustGet(common.PluginPubSub).(core.PubSubComponent)

		repo := mysql.NewMySQLRepository(db)
		business := business.NewUserLikeItemBusiness(repo, ps)
		now := time.Now().UTC()

		if err := business.LikeItem(ctx.Request.Context(), &entity.Like{
			UserID:    common.GetRequesterID(requester),
			ItemID:    int(id.GetLocalID()),
			CreatedAt: &now,
		}); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(true))
	}
}
