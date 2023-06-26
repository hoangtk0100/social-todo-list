package ginuserlikeitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
)

func ListLikedUsers(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := common.UIDFromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var paging common.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Process()

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		store := storage.NewSQLStore(db)
		business := biz.NewListUsersLikedItemBiz(store)

		result, err := business.ListUsersLikedItem(ctx.Request.Context(), int(id.GetLocalID()), &paging)
		if err != nil {
			panic(err)
		}

		for index := range result {
			result[index].Mask()
		}

		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))
	}
}
