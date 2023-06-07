package ginitem

import (
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/biz"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
	"gorm.io/gorm"
)

func ListItem(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
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
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store, requester)

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
