package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/component/uploadprovider"
	"github.com/hoangtk0100/social-todo-list/middleware"
	ginitem "github.com/hoangtk0100/social-todo-list/module/item/transport/gin"
	"github.com/hoangtk0100/social-todo-list/module/upload/transport/ginupload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_SOURCE")
	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	storageSecretKey := os.Getenv("STORAGE_SECRET_KEY")
	storageRegion := os.Getenv("STORAGE_REGION")
	storageBucket := os.Getenv("STORAGE_BUCKET")
	storageEndPoint := os.Getenv("STORAGE_END_POINT")
	storageDomain := os.Getenv("STORAGE_DOMAIN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	log.Println("DB Connected :", db)

	r2Provider := uploadprovider.NewR2Provider(storageBucket, storageRegion, storageAccessKey, storageSecretKey, storageEndPoint, storageDomain)

	r := gin.Default()
	r.Use(middleware.Recover())

	r.Static("/static", "./static")

	v1 := r.Group("/v1")
	{
		uploads := v1.Group("/upload")
		{
			uploads.POST("", ginupload.Upload(db, r2Provider))
			uploads.POST("/local", ginupload.UploadLocal(db))
		}

		items := v1.Group("/items")
		{
			items.POST("", ginitem.CreateItem(db))
			items.GET("", ginitem.ListItem(db))
			items.GET("/:id", ginitem.GetItem(db))
			items.PATCH("/:id", ginitem.UpdateItem(db))
			items.DELETE("/:id", ginitem.DeleteItem(db))
		}
	}

	r.Run(":9090")
}
