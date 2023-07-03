package api

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/business"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
)

func GetItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(entity.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		repo := mysql.NewMySQLRepository(db)
		business := business.NewGetItemBusiness(repo)

		data, err := business.GetItemByID(ctx.Request.Context(), int(id.GetLocalID()))
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		data.Mask()

		core.SuccessResponse(ctx, core.NewDataResponse(data))
	}
}
