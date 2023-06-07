package ginupload

import (
	"mime/multipart"
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/upload/biz"
	"github.com/hoangtk0100/social-todo-list/module/upload/storage"
	"github.com/hoangtk0100/social-todo-list/plugin/uploadprovider"
	"gorm.io/gorm"
)

func Upload(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		_, dataBytes, folder, fileName, contentType := validateFiles(ctx)

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		provider := serviceCtx.MustGet(common.PluginR2).(uploadprovider.UploadProvider)
		store := storage.NewSQLStore(db)
		business := biz.NewUploadBiz(store, provider)

		img, err := business.Upload(ctx.Request.Context(), dataBytes, folder, fileName, contentType)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}

func UploadLocal() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileHeader, dataBytes, folder, fileName, _ := validateFiles(ctx)

		business := biz.NewUploadBiz(nil, nil)
		img, err := business.UploadLocal(ctx, fileHeader, dataBytes, folder, fileName)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}

func validateFiles(ctx *gin.Context) (fileHeader *multipart.FileHeader, dataBytes []byte, folder, fileName, contentType string) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		panic(common.ErrInvalidRequest(err))
	}

	folder = ctx.DefaultPostForm("folder", "images")
	file, err := fileHeader.Open()
	if err != nil {
		panic(common.ErrInvalidRequest(err))
	}

	defer file.Close()

	fileName = fileHeader.Filename
	contentType = fileHeader.Header.Get("Content-Type")
	dataBytes = make([]byte, fileHeader.Size)
	if _, err := file.Read(dataBytes); err != nil {
		panic(common.ErrInvalidRequest(err))
	}

	return
}
