package biz

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/upload/model"
)

type CreateImageStorage interface {
	CreateImage(ctx context.Context, data *common.Image) error
}

type uploadBiz struct {
	store    CreateImageStorage
	provider core.StorageComponent
}

func NewUploadBiz(store CreateImageStorage, provider core.StorageComponent) *uploadBiz {
	return &uploadBiz{store: store, provider: provider}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string, contentType string) (*common.Image, error) {
	fileBytes := bytes.NewBuffer(data)

	width, height, err := getImageDimension(fileBytes)
	if err != nil {
		return nil, model.ErrFileNotImage(err)
	}

	fileExt := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("%d.%s", time.Now().UTC().UnixNano(), fileName)
	dst := fmt.Sprintf("%s/%s", folder, newFileName)
	url, storageName, err := biz.provider.UploadFile(ctx, data, dst, contentType)
	if err != nil {
		return nil, model.ErrCannotSaveFile(err)
	}

	img := &common.Image{
		Name:      fileName,
		Width:     width,
		Height:    height,
		Extension: getShortExtension(fileExt),
		Url:       url,
		CloudName: storageName,
	}

	if err := biz.store.CreateImage(ctx, img); err != nil {
		biz.provider.DeleteFiles(ctx, []string{dst})
		return nil, model.ErrCannotSaveFile(err)
	}

	return img, nil
}

func (biz *uploadBiz) UploadLocal(ctx context.Context, fileHeader *multipart.FileHeader, data []byte, folder, fileName string) (*common.Image, error) {
	fileBytes := bytes.NewBuffer(data)

	width, height, err := getImageDimension(fileBytes)
	if err != nil {
		return nil, model.ErrFileNotImage(err)
	}

	fileExt := filepath.Ext(fileName)
	dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().UnixNano(), fileName)

	if err := ctx.(*gin.Context).SaveUploadedFile(fileHeader, dst); err != nil {
		return nil, model.ErrCannotSaveFile(err)
	}

	img := &common.Image{
		Url:       dst,
		Width:     width,
		Height:    height,
		CloudName: "local",
		Extension: getShortExtension(fileExt),
	}

	img.Fulfill("http://localhost:9090")
	return img, nil
}

func getShortExtension(ext string) string {
	return strings.ReplaceAll(ext, ".", "")
}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		log.Println("err:", err)
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}
