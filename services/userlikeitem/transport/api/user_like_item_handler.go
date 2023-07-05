package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

func (service *service) LikeItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.UIDFromString(ctx.Param("id"))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.
				WithError(entity.ErrItemIDInvalid.Error()).
				WithDebug(err.Error()),
			)

			return
		}

		requester := core.GetRequester(ctx)
		now := time.Now().UTC()

		if err := service.business.LikeItem(ctx, &entity.Like{
			UserID:    common.GetRequesterID(requester),
			ItemID:    int(id.GetLocalID()),
			CreatedAt: &now,
		}); err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(true))
	}
}
