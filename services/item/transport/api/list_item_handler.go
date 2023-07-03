package api

import (
	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/business"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/rpc"
)

func ListItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var queryString struct {
			core.Paging
			entity.Filter
		}

		if err := ctx.ShouldBind(&queryString); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		queryString.Paging.Process()

		requester := core.GetRequester(ctx)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		itemAPICaller := ac.MustGet(common.PluginItemAPI).(interface {
			GetServiceURL() string
		})

		repo := mysql.NewMySQLRepository(db)
		likeRepo := rpc.NewItemAPIClient(itemAPICaller.GetServiceURL(), ac.Logger("rpc.itemlikes"))
		business := business.NewListItemBusiness(repo, likeRepo, requester)

		result, err := business.ListItem(ctx.Request.Context(), &queryString.Filter, &queryString.Paging)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		for index := range result {
			result[index].Mask()
		}

		core.SuccessResponse(ctx, core.NewResponse(result, queryString.Paging, queryString.Filter))
	}
}
