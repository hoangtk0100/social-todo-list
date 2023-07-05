package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (service *service) GetItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(entity.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		data, err := service.business.GetItemByID(ctx, int(id.GetLocalID()))
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		data.Mask()
		core.SuccessResponse(ctx, core.NewDataResponse(data))
	}
}
