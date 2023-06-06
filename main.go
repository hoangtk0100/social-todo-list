package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/component/uploadprovider"
	"github.com/hoangtk0100/social-todo-list/middleware"
	ginitem "github.com/hoangtk0100/social-todo-list/module/item/transport/gin"
	ginupload "github.com/hoangtk0100/social-todo-list/module/upload/transport/gin"
	ginuser "github.com/hoangtk0100/social-todo-list/module/user/transport/gin"
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
	serverAddress := os.Getenv("SERVER_ADDRESS")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	log.Println("DB Connected :", db)

	r2Provider := uploadprovider.NewR2Provider(storageBucket, storageRegion, storageAccessKey, storageSecretKey, storageEndPoint, storageDomain)

	router := setupRoutes(db, r2Provider)
	router.Run(fmt.Sprint(":", serverAddress))
}

func setupRoutes(db *gorm.DB, uploadProvider uploadprovider.UploadProvider) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Recover())

	router.Static("/static", "./static")

	v1 := router.Group("/v1")
	{
		v1.POST("/register", ginuser.Register(db))

		uploads := v1.Group("/upload")
		{
			uploads.POST("", ginupload.Upload(db, uploadProvider))
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

	return router
}
