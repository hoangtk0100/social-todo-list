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

func UpdateItem(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.TodoItemUpdate
		if err := ctx.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewSQLStore(db)
		business := biz.NewUpdateItemBiz(store, requester)

		if err := business.UpdateItemById(ctx.Request.Context(), int(id.GetLocalID()), &data); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
