package rpc

import (
	"context"

	"fmt"

	"github.com/go-resty/resty/v2"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/pkg/errors"
)

type itemService struct {
	client     *resty.Client
	serviceURL string
	logger     appctx.Logger
}

func NewItemService(serviceURL string, logger appctx.Logger) *itemService {
	return &itemService{
		client:     resty.New(),
		serviceURL: serviceURL,
		logger:     logger,
	}
}

func (s *itemService) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	type requestBody struct {
		Ids []int `json:"ids"`
	}

	var response struct {
		Data map[int]int `json:"data"`
	}

	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody{Ids: ids}).
		SetResult(&response).
		Post(fmt.Sprintf("%s/%s", s.serviceURL, "v1/rpc/get_item_likes"))

	if err != nil {
		s.logger.Error(err)
		return nil, errors.WithStack(err)
	}

	if !resp.IsSuccess() {
		respErr := resp.Error().(error)
		s.logger.Errorf(respErr, resp.String())
		return nil, errors.WithStack(respErr)
	}

	return response.Data, nil
}
