package ginupload

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/upload/biz"
	"github.com/hoangtk0100/social-todo-list/module/upload/storage"
)

func Upload(ac appctx.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, dataBytes, folder, fileName, contentType := validateFiles(ctx)

		db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		provider := ac.MustGet(common.PluginR2).(core.StorageComponent)
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
