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
	itemBusiness "github.com/hoangtk0100/social-todo-list/services/item/business"
	itemSQLRepo "github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/rpc"
	itemAPI "github.com/hoangtk0100/social-todo-list/services/item/transport/api"
	uploadBusiness "github.com/hoangtk0100/social-todo-list/services/upload/business"
	uploadSQLRepo "github.com/hoangtk0100/social-todo-list/services/upload/repository/mysql"
	uploadAPI "github.com/hoangtk0100/social-todo-list/services/upload/transport/api"
	userBusiness "github.com/hoangtk0100/social-todo-list/services/user/business"
	userSQLRepo "github.com/hoangtk0100/social-todo-list/services/user/repository/mysql"
	userAPI "github.com/hoangtk0100/social-todo-list/services/user/transport/api"
	userLikeItemBusiness "github.com/hoangtk0100/social-todo-list/services/userlikeitem/business"
	userLikeItemSQLRepo "github.com/hoangtk0100/social-todo-list/services/userlikeitem/repository/mysql"
	userLikeItemAPI "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/api"
	userlikeitemRPC "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/rpc"
	"github.com/hoangtk0100/social-todo-list/subscriber"
	"github.com/spf13/cobra"
)

func newAppContext() appctx.AppContext {
	return appctx.NewAppContext(
		appctx.WithName("social-todo-list"),
		appctx.WithComponent(gormdb.NewGormDB(common.PluginDBMain, common.PluginDBMain)),
		appctx.WithComponent(token.NewJWTMaker(common.PluginJWT)),
		appctx.WithComponent(appstorage.NewR2Storage(common.PluginR2)),
		appctx.WithComponent(tracer.NewJaeger(common.PluginTracerJaeger)),
		appctx.WithComponent(pubsub.NewNatsPubSub(common.PluginPubSub)),
		appctx.WithComponent(redisdb.NewRedisDB(common.PluginRedis, common.PluginRedis)),
		appctx.WithComponent(rpccaller.NewItemAPICaller(common.PluginItemAPI)),
		appctx.WithComponent(ginserver.NewGinServer(common.PluginGin)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		appCtx := newAppContext()
		log := appCtx.Logger("service")

		if err := appCtx.Load(); err != nil {
			log.Fatal(err)
		}

		ginServer := appCtx.MustGet(common.PluginGin).(core.GinComponent)
		router := ginServer.GetRouter()
		router.Use(ginmiddleware.Recovery(appCtx))

		db := appCtx.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
		tokenMaker := appCtx.MustGet(common.PluginJWT).(core.TokenMakerComponent)
		ps := appCtx.MustGet(common.PluginPubSub).(core.PubSubComponent)
		storageProvider := appCtx.MustGet(common.PluginR2).(core.StorageComponent)

		userRepo := userSQLRepo.NewMySQLRepository(db)
		userBiz := userBusiness.NewBusiness(userRepo, tokenMaker)
		userService := userAPI.NewService(appCtx, userBiz)

		authCache := memcache.NewUserCache(cache.NewRedisCache(common.PluginRedis, appCtx), userRepo)
		authMiddleware := middleware.RequireAuth(authCache, appCtx)

		itemAPICaller := appCtx.MustGet(common.PluginItemAPI).(common.ItemAPICaller)
		likeRepo := rpc.NewItemAPIClient(itemAPICaller.GetServiceURL(), appCtx.Logger("rpc.itemlikes"))

		itemRepo := itemSQLRepo.NewMySQLRepository(db)
		itemBiz := itemBusiness.NewBusiness(itemRepo, likeRepo)
		itemService := itemAPI.NewService(appCtx, itemBiz)

		userLikeItemRepo := userLikeItemSQLRepo.NewMySQLRepository(db)
		userLikeItemBiz := userLikeItemBusiness.NewBusiness(userLikeItemRepo, ps)
		userLikeItemService := userLikeItemAPI.NewService(appCtx, userLikeItemBiz)

		uploadRepo := uploadSQLRepo.NewMySQLRepository(db)
		uploadBiz := uploadBusiness.NewBusiness(uploadRepo, storageProvider)
		uploadService := uploadAPI.NewService(appCtx, uploadBiz)

		router.Static("/static", "./static")
		v1 := router.Group("/v1")
		{
			v1.POST("/register", userService.Register())
			v1.POST("/login", userService.Login())
			v1.GET("/profile", authMiddleware, userService.Profile())

			uploads := v1.Group("/upload", authMiddleware)
			{
				uploads.POST("", uploadService.Upload())
				uploads.POST("/local", uploadService.UploadLocal())
			}

			items := v1.Group("/items", authMiddleware)
			{
				items.POST("", itemService.CreateItem())
				items.GET("", itemService.ListItem())
				items.GET("/:id", itemService.GetItem())
				items.PATCH("/:id", itemService.UpdateItem())
				items.DELETE("/:id", itemService.DeleteItem())

				items.POST("/:id/like", userLikeItemService.LikeItem())
				items.DELETE("/:id/unlike", userLikeItemService.UnlikeItem())
				items.GET("/:id/liked-users", userLikeItemService.ListLikedUsers())
			}

			rpc := v1.Group("/rpc")
			{
				rpc.POST("/get_item_likes", userlikeitemRPC.GetItemLikes(appCtx))
			}
		}

		subscriber.StartPbEngine(appCtx)

		ginServer.Start()
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
