package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		user := ctx.MustGet(common.CurrentUser)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
