package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/biz"
	"github.com/hoangtk0100/social-todo-list/module/item/model"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
	"gorm.io/gorm"
)

func ListItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		queryString.Paging.Process()

		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
