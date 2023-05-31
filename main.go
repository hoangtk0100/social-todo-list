package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	ginitem "github.com/hoangtk0100/social-todo-list/module/item/transport/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_SOURCE")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	log.Println("DB Connected :", db)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", ginitem.CreateItem(db))
		}
	}

	r.Run(":9090")
}
