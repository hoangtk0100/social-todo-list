package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
)

func (service *service) Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := service.business.GetProfile(ctx)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		data := core.NewSimpleUser(user.ID, user.FirstName, user.LastName, user.Avatar)
		data.Mask(common.MaskTypeUser)
		core.SuccessResponse(ctx, core.NewDataResponse(data))
	}
}
