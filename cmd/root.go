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
	ginmiddleware "github.com/hoangtk0100/app-context/component/server/gin/middleware"
	appstorage "github.com/hoangtk0100/app-context/component/storage"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/component/tracer"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/rpccaller"
	"github.com/hoangtk0100/social-todo-list/memcache"
	"github.com/hoangtk0100/social-todo-list/middleware"
	itemAPI "github.com/hoangtk0100/social-todo-list/services/item/transport/api"
	uploadAPI "github.com/hoangtk0100/social-todo-list/services/upload/transport/api"
	userRepo "github.com/hoangtk0100/social-todo-list/services/user/repository/mysql"
	userAPI "github.com/hoangtk0100/social-todo-list/services/user/transport/api"
	userlikeitemAPI "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/api"
	userlikeitemRPC "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/rpc"
	"github.com/hoangtk0100/social-todo-list/subscriber"
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
		appctx.WithComponent(rpccaller.NewitemAPICaller(common.PluginItemAPI)),
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
		router.Use(ginmiddleware.Recovery(service))

		db := service.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		authRepo := userRepo.NewMySQLRepository(db)
		authCache := memcache.NewUserCache(cache.NewRedisCache(common.PluginRedis, service), authRepo)
		authMiddleware := middleware.RequireAuth(authCache, service)

		router.Static("/static", "./static")
		v1 := router.Group("/v1")
		{
			v1.POST("/register", userAPI.Register(service))
			v1.POST("/login", userAPI.Login(service))
			v1.GET("/profile", authMiddleware, userAPI.Profile(service))

			uploads := v1.Group("/upload", authMiddleware)
			{
				uploads.POST("", uploadAPI.Upload(service))
				uploads.POST("/local", uploadAPI.UploadLocal())
			}

			items := v1.Group("/items", authMiddleware)
			{
				items.POST("", itemAPI.CreateItem(service))
				items.GET("", itemAPI.ListItem(service))
				items.GET("/:id", itemAPI.GetItem(service))
				items.PATCH("/:id", itemAPI.UpdateItem(service))
				items.DELETE("/:id", itemAPI.DeleteItem(service))

				items.POST("/:id/like", userlikeitemAPI.LikeItem(service))
				items.DELETE("/:id/unlike", userlikeitemAPI.UnlikeItem(service))
				items.GET("/:id/liked-users", userlikeitemAPI.ListLikedUsers(service))
			}

			rpc := v1.Group("/rpc")
			{
				rpc.POST("/get_item_likes", userlikeitemRPC.GetItemLikes(service))
			}
		}

		subscriber.StartPbEngine(service)
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
