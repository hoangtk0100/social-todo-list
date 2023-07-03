package biz

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/module/upload/model"
)

type CreateImageStorage interface {
	CreateImage(ctx context.Context, data *core.Image) error
}

type uploadBiz struct {
	store    CreateImageStorage
	provider core.StorageComponent
}

func NewUploadBiz(store CreateImageStorage, provider core.StorageComponent) *uploadBiz {
	return &uploadBiz{store: store, provider: provider}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string, contentType string) (*core.Image, error) {
	width, height, err := util.GetImageDimension(data)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(model.ErrFileNotImage.Error()).WithDebug(err.Error())
	}

	newFileName := fmt.Sprintf("%d.%s", time.Now().UTC().UnixNano(), fileName)
	dst := fmt.Sprintf("%s/%s", folder, newFileName)
	url, storageName, err := biz.provider.UploadFile(ctx, data, dst, contentType)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(model.ErrCannotSaveFile.Error()).WithDebug(err.Error())
	}

	img := &core.Image{
		Name:     fileName,
		Path:     dst,
		Width:    width,
		Height:   height,
		URL:      url,
		Provider: storageName,
	}

	if err := biz.store.CreateImage(ctx, img); err != nil {
		biz.provider.DeleteFiles(ctx, []string{dst})
		return nil, core.ErrBadRequest.WithError(model.ErrCannotSaveFile.Error()).WithDebug(err.Error())
	}

	return img, nil
}

func (biz *uploadBiz) UploadLocal(ctx context.Context, fileHeader *multipart.FileHeader, data []byte, folder, fileName string) (*core.Image, error) {
	width, height, err := util.GetImageDimension(data)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(model.ErrFileNotImage.Error()).WithDebug(err.Error())
	}

	dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().UnixNano(), fileName)

	if err := ctx.(*gin.Context).SaveUploadedFile(fileHeader, dst); err != nil {
		return nil, core.ErrBadRequest.WithError(model.ErrCannotSaveFile.Error()).WithDebug(err.Error())
	}

	img := &core.Image{
		Name:     fileName,
		Path:     dst,
		Width:    width,
		Height:   height,
		Provider: "local",
	}

	img.Fulfill("http://localhost:9090")
	return img, nil
}
