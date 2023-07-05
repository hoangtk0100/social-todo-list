package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/entity"
)

func (service *service) ListLikedUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(entity.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		var paging core.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		paging.Process()

		result, err := service.business.ListUsersLikedItem(ctx, int(id.GetLocalID()), &paging)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		for index := range result {
			result[index].Mask(common.MaskTypeUser)
		}

		core.SuccessResponse(ctx, core.NewResponse(result, paging, nil))
	}
}
