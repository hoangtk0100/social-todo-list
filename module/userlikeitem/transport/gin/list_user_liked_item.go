package ginuserlikeitem

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/model"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func ListLikedUsers(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(model.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		var paging core.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		paging.Process()

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		store := storage.NewSQLStore(db)
		business := biz.NewListUsersLikedItemBiz(store)

		result, err := business.ListUsersLikedItem(ctx.Request.Context(), int(id.GetLocalID()), &paging)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		for index := range result {
			result[index].Mask(common.MaskTypeUser)
		}

		core.SuccessResponse(ctx, core.NewResponse(result, paging, nil))
	}
}
