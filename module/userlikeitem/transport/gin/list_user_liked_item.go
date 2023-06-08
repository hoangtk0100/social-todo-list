package ginuserlikeitem

import (
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/biz"
	"github.com/hoangtk0100/social-todo-list/module/userlikeitem/storage"
	"gorm.io/gorm"
)

func ListLikedUsers(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var paging common.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Process()

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

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
