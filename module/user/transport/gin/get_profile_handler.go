package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

func Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet(common.CurrentUser)

		user.(*model.User).SQLModel.Mask(common.DBTypeUser)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
