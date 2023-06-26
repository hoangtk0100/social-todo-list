package cmd

import (
	"fmt"
	"os"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/cache"
	"github.com/hoangtk0100/app-context/component/datastore/gormdb"
	"github.com/hoangtk0100/app-context/component/datastore/redisdb"
	"github.com/hoangtk0100/app-context/component/pubsub"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
	appstorage "github.com/hoangtk0100/app-context/component/storage"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/component/tracer"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/memcache"
	"github.com/hoangtk0100/social-todo-list/middleware"
	ginitem "github.com/hoangtk0100/social-todo-list/module/item/transport/gin"
	ginuserlikeitem "github.com/hoangtk0100/social-todo-list/module/userlikeitem/transport/gin"
	rpcuserlikeitem "github.com/hoangtk0100/social-todo-list/module/userlikeitem/transport/rpc"
	"github.com/hoangtk0100/social-todo-list/subscriber"

	ginupload "github.com/hoangtk0100/social-todo-list/module/upload/transport/gin"
	userstorage "github.com/hoangtk0100/social-todo-list/module/user/storage"
	ginuser "github.com/hoangtk0100/social-todo-list/module/user/transport/gin"
	"github.com/hoangtk0100/social-todo-list/plugin/rpccaller"

	"github.com/spf13/cobra"
)

func newService() appctx.AppContext {
	return appctx.NewAppContext(
		appctx.WithName("social-todo-list"),
		appctx.WithComponent(gormdb.NewGormDB(common.PluginDBMain, common.PluginDBMain)),
		appctx.WithComponent(token.NewJWTMaker(common.PluginJWT)),
		appctx.WithComponent(appstorage.NewR2Storage(common.PluginR2)),
		appctx.WithComponent(tracer.NewJaeger(common.PluginTracerJaeger)),
		appctx.WithComponent(pubsub.NewNatsPubSub(common.PluginPubSub)),
		appctx.WithComponent(redisdb.NewRedisDB(common.PluginRedis, common.PluginRedis)),
		appctx.WithComponent(rpccaller.NewApiItemCaller(common.PluginItemAPI)),
		appctx.WithComponent(ginserver.NewGinServer(common.PluginGin)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()
		log := service.Logger("service")

		if err := service.Load(); err != nil {
			log.Fatal(err)
		}

		server := service.MustGet(common.PluginGin).(core.GinComponent)
		router := server.GetRouter()
		router.Use(middleware.Recover())

		db := service.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		authStore := userstorage.NewSQLStore(db)
		authCache := memcache.NewUserCache(cache.NewRedisCache(common.PluginRedis, service), authStore)
		authMiddleware := middleware.RequireAuth(authCache, service)

		router.Static("/static", "./static")
		v1 := router.Group("/v1")
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

			rpc := v1.Group("/rpc")
			{
				rpc.POST("/get_item_likes", rpcuserlikeitem.GetItemLikes(service))
			}
		}

		_ = subscriber.NewPBEngine(service).Start()

		server.Start()
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
