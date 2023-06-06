package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/tokenprovider"
	"github.com/hoangtk0100/social-todo-list/module/user/biz"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/module/user/storage"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, tokenProvider tokenprovider.TokenProvider) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var loginData model.UserLogin

		if err := ctx.ShouldBind(&loginData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

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
