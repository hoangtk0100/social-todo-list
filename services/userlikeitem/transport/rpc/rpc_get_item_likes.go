package rpc

import (
	"context"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/userlikeitem/entity"
)

type UserLikeItemBusiness interface {
	LikeItem(ctx context.Context, data *entity.Like) error
	UnlikeItem(ctx context.Context, userID, itemID int) error
	ListUsersLikedItem(ctx context.Context, itemID int, paging *core.Paging) ([]core.SimpleUser, error)
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type service struct {
	ac       appctx.AppContext
	business UserLikeItemBusiness
}

func NewService(ac appctx.AppContext, business UserLikeItemBusiness) *service {
	return &service{
		ac:       ac,
		business: business,
	}
}

func (service *service) GetItemLikes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type requestData struct {
			IDs []int `json:"ids"`
		}

		var data requestData
		if err := ctx.ShouldBind(&data); err != nil {
			core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		mapResult, err := service.business.GetItemLikes(ctx.Request.Context(), data.IDs)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(mapResult))
	}
}
