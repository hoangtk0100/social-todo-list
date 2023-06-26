package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/biz"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/hoangtk0100/social-todo-list/module/item/repository"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
	"github.com/hoangtk0100/social-todo-list/module/item/storage/rpc"
)

func ListItem(ac appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := ctx.ShouldBind(&queryString); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		queryString.Paging.Process()

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		apiItemCaller := ac.MustGet(common.PluginItemAPI).(interface {
			GetServiceURL() string
		})

		store := storage.NewSQLStore(db)
		likeStore := rpc.NewItemService(apiItemCaller.GetServiceURL(), ac.Logger("rpc.itemlikes"))
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)

		result, err := business.ListItem(ctx.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			panic(err)
		}

		for index := range result {
			result[index].Mask()
		}

		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
