package rpcuserlikeitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func GetItemLikes(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type requestData struct {
			Ids []int `json:"ids"`
		}

		var data requestData
		if err := ctx.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		store := storage.NewSQLStore(db)

		mapResult, err := store.GetItemLikes(ctx.Request.Context(), data.Ids)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(mapResult))
	}
}
