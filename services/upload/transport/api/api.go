package api

import (
	"context"
	"mime/multipart"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
)

type UploadBusiness interface {
	Upload(ctx context.Context, data []byte, folder, fileName string, contentType string) (*core.Image, error)
	UploadLocal(ctx context.Context, fileHeader *multipart.FileHeader, data []byte, folder, fileName string) (*core.Image, error)
}

type service struct {
	ac       appctx.AppContext
	business UploadBusiness
}

func NewService(ac appctx.AppContext, business UploadBusiness) *service {
	return &service{
		ac:       ac,
		business: business,
	}
}
