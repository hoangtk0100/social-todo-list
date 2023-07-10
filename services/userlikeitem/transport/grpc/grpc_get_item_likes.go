package grpc

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/pb"
)

type UserLikeItemBusiness interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type grpcService struct {
	business UserLikeItemBusiness
}

func NewService(business UserLikeItemBusiness) *grpcService {
	return &grpcService{
		business: business,
	}
}

func (service *grpcService) GetItemLikes(ctx context.Context, req *pb.GetItemLikesRequest) (*pb.GetItemLikesResponse, error) {
	ids := util.ConvertInt32SliceToIntSlice(req.Ids)
	mapResult, err := service.business.GetItemLikes(ctx, ids)
	if err != nil {
		return nil, core.ErrInternalServerError.WithError(err.Error())
	}

	data := util.ConvertIntMapToInt32Map(mapResult)
	return &pb.GetItemLikesResponse{Data: data}, nil
}
