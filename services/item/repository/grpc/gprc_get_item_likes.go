package grpc

import (
	"context"

	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/pb"
	"github.com/pkg/errors"
)

type grpcClient struct {
	client pb.UserLikeItemServiceClient
}

func NewClient(client pb.UserLikeItemServiceClient) *grpcClient {
	return &grpcClient{client: client}
}

func (c *grpcClient) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	reqIDs := util.ConvertIntSliceToInt32Slice(ids)
	resp, err := c.client.GetItemLikes(ctx, &pb.GetItemLikesRequest{Ids: reqIDs})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return util.ConvertInt32MapToIntMap(resp.Data), nil
}
