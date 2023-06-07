package ginuser

import (
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/biz"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/user/storage"
	"github.com/hoangtk0100/social-todo-list/plugin/tokenprovider"
	"gorm.io/gorm"
)

func Login(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var loginData model.UserLogin

		if err := ctx.ShouldBind(&loginData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.TokenProvider)

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		expiry := 60 * 60 * 24 * 7
		business := biz.NewLoginBiz(store, tokenProvider, md5, expiry)
		token, err := business.Login(ctx.Request.Context(), &loginData)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(token))
	}
}
