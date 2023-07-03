package api

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/user/business"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
	"github.com/hoangtk0100/social-todo-list/services/user/repository/mysql"
)

func Login(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData entity.UserLogin

		if err := ctx.ShouldBind(&loginData); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		tokenMaker := ac.MustGet(common.PluginJWT).(core.TokenMakerComponent)

		repo := mysql.NewMySQLRepository(db)
		business := business.NewLoginBusiness(repo, tokenMaker)
		tokenPayload, err := business.Login(ctx.Request.Context(), &loginData)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(tokenPayload))
	}
}
