package ginupload

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/upload/biz"
	"github.com/hoangtk0100/social-todo-list/module/upload/model"
	"github.com/hoangtk0100/social-todo-list/module/upload/storage"
)

func Upload(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, dataBytes, folder, fileName, contentType, err := validateFiles(ctx)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		provider := ac.MustGet(common.PluginR2).(core.StorageComponent)
		store := storage.NewSQLStore(db)
		business := biz.NewUploadBiz(store, provider)

		img, err := business.Upload(ctx.Request.Context(), dataBytes, folder, fileName, contentType)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(img))
	}
}

func UploadLocal() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileHeader, dataBytes, folder, fileName, _, err := validateFiles(ctx)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		business := biz.NewUploadBiz(nil, nil)
		img, err := business.UploadLocal(ctx, fileHeader, dataBytes, folder, fileName)
		if err != nil {
			core.ErrorResponse(ctx, err)
			return
		}

		core.SuccessResponse(ctx, core.NewDataResponse(img))
	}
}

func validateFiles(ctx *gin.Context) (fileHeader *multipart.FileHeader, dataBytes []byte, folder, fileName, contentType string, err error) {
	fileHeader, err = ctx.FormFile("file")
	if err != nil {
		err = core.ErrBadRequest.
			WithError(model.ErrFileMissing.Error()).
			WithDebug(err.Error())
		return
	}

	folder = ctx.DefaultPostForm("folder", "images")
	file, err := fileHeader.Open()
	if err != nil {
		err = core.ErrBadRequest.
			WithError(model.ErrCannotReadFile.Error()).
			WithDebug(err.Error())
		return
	}

	defer file.Close()

	fileName = fileHeader.Filename
	contentType = fileHeader.Header.Get("Content-Type")
	dataBytes = make([]byte, fileHeader.Size)
	if _, err = file.Read(dataBytes); err != nil {
		err = core.ErrBadRequest.
			WithError(model.ErrCannotReadFile.Error()).
			WithDebug(err.Error())
		return
	}

	return
}
