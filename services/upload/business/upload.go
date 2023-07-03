package business

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/social-todo-list/services/upload/entity"
)

type CreateImageRepository interface {
	CreateImage(ctx context.Context, data *core.Image) error
}

type uploadBusiness struct {
	repo     CreateImageRepository
	provider core.StorageComponent
}

func NewUploadBusiness(repo CreateImageRepository, provider core.StorageComponent) *uploadBusiness {
	return &uploadBusiness{repo: repo, provider: provider}
}

func (biz *uploadBusiness) Upload(ctx context.Context, data []byte, folder, fileName string, contentType string) (*core.Image, error) {
	width, height, err := util.GetImageDimension(data)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(entity.ErrFileNotImage.Error()).WithDebug(err.Error())
	}

	newFileName := fmt.Sprintf("%d.%s", time.Now().UTC().UnixNano(), fileName)
	dst := fmt.Sprintf("%s/%s", folder, newFileName)
	url, storageName, err := biz.provider.UploadFile(ctx, data, dst, contentType)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(entity.ErrCannotSaveFile.Error()).WithDebug(err.Error())
	}

	img := &core.Image{
		Name:     fileName,
		Path:     dst,
		Width:    width,
		Height:   height,
		URL:      url,
		Provider: storageName,
	}

	if err := biz.repo.CreateImage(ctx, img); err != nil {
		biz.provider.DeleteFiles(ctx, []string{dst})
		return nil, core.ErrBadRequest.WithError(entity.ErrCannotSaveFile.Error()).WithDebug(err.Error())
	}

	return img, nil
}

func (biz *uploadBusiness) UploadLocal(ctx context.Context, fileHeader *multipart.FileHeader, data []byte, folder, fileName string) (*core.Image, error) {
	width, height, err := util.GetImageDimension(data)
	if err != nil {
		return nil, core.ErrBadRequest.WithError(entity.ErrFileNotImage.Error()).WithDebug(err.Error())
	}

	dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().UnixNano(), fileName)

	if err := ctx.(*gin.Context).SaveUploadedFile(fileHeader, dst); err != nil {
		return nil, core.ErrBadRequest.WithError(entity.ErrCannotSaveFile.Error()).WithDebug(err.Error())
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
