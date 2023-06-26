package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	appctx "github.com/hoangtk0100/app-context"
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
		return nil, err
	}

	if !resp.IsSuccess() {
		s.logger.Errorf(resp.Error().(error), resp.String())
		return nil, errors.New("cannot call get item likes")
	}

	return response.Data, nil
}
