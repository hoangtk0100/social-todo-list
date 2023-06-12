package cmd

import (
	"fmt"
	"os"

	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/middleware"
	ginitem "github.com/hoangtk0100/social-todo-list/module/item/transport/gin"
	ginuserlikeitem "github.com/hoangtk0100/social-todo-list/module/userlikeitem/transport/gin"
	pubsub "github.com/hoangtk0100/social-todo-list/pubsub"
	"github.com/hoangtk0100/social-todo-list/subscriber"

	ginupload "github.com/hoangtk0100/social-todo-list/module/upload/transport/gin"
	userstorage "github.com/hoangtk0100/social-todo-list/module/user/storage"
	ginuser "github.com/hoangtk0100/social-todo-list/module/user/transport/gin"
	"github.com/hoangtk0100/social-todo-list/plugin/sdkgorm"
	"github.com/hoangtk0100/social-todo-list/plugin/tokenprovider/jwt"
	"github.com/hoangtk0100/social-todo-list/plugin/uploadprovider"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		goservice.WithInitRunnable(uploadprovider.NewR2Provider(common.PluginR2)),
		goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()
		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			authMiddleware := middleware.RequireAuth(authStore, service)

			engine.Static("/static", "./static")
			v1 := engine.Group("/v1")
			{
				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.GET("/profile", authMiddleware, ginuser.Profile())

				uploads := v1.Group("/upload", authMiddleware)
				{
					uploads.POST("", ginupload.Upload(service))
					uploads.POST("/local", ginupload.UploadLocal())
				}

				items := v1.Group("/items", authMiddleware)
				{
					items.POST("", ginitem.CreateItem(service))
					items.GET("", ginitem.ListItem(service))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))

					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListLikedUsers(service))
				}
			}
		})

		_ = subscriber.NewPBEngine(service).Start()

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronUpdateItemLikedCountCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
