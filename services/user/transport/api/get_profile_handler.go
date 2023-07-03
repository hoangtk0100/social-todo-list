package api

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/business"
	"github.com/hoangtk0100/social-todo-list/services/user/repository/mysql"
)

func Profile(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		repo := mysql.NewMySQLRepository(db)
		business := business.NewGetProfileBusiness(repo)

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
