package upload

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().UnixNano(), fileHeader.Filename)

		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {

		}

		img := common.Image{
			Url:       dst,
			Width:     100,
			Height:    100,
			CloudName: "local",
			Extension: "",
		}

		img.Fulfill("http://localhost:9090")

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
