package ginupload

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/uploadprovider"
	"github.com/hoangtk0100/social-todo-list/module/upload/biz"
	"github.com/hoangtk0100/social-todo-list/module/upload/storage"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB, provider uploadprovider.UploadProvider) func(*gin.Context) {
	return func(c *gin.Context) {
		_, dataBytes, folder, fileName, contentType := validateFiles(c)
		store := storage.NewSQLStore(db)
		business := biz.NewUploadBiz(store, provider)

		img, err := business.Upload(c.Request.Context(), dataBytes, folder, fileName, contentType)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}

func UploadLocal(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, dataBytes, folder, fileName, _ := validateFiles(c)

		business := biz.NewUploadBiz(nil, nil)
		img, err := business.UploadLocal(c, fileHeader, dataBytes, folder, fileName)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
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
