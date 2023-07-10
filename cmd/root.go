package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hoangtk0100/social-todo-list/pb"

	"github.com/gin-gonic/gin"
	appctx "github.com/hoangtk0100/app-context"
	grpcclient "github.com/hoangtk0100/app-context/component/client/grpc"
	"github.com/hoangtk0100/app-context/component/datastore/gormdb"
	"github.com/hoangtk0100/app-context/component/datastore/redisdb"
	"github.com/hoangtk0100/app-context/component/pubsub"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
	"github.com/hoangtk0100/app-context/component/server/gin/middleware"
	grpcserver "github.com/hoangtk0100/app-context/component/server/grpc"
	"github.com/hoangtk0100/app-context/component/storage"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/component/tracer"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/builder"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/rpccaller"
	"github.com/hoangtk0100/social-todo-list/subscriber"
	"github.com/spf13/cobra"
)

func newAppContext() appctx.AppContext {
	return appctx.NewAppContext(
		appctx.WithName("social-todo-list"),
		appctx.WithComponent(gormdb.NewGormDB(common.PluginDBMain, common.PluginDBMain)),
		appctx.WithComponent(token.NewPasetoMaker(common.PluginTokenMaker)),
		appctx.WithComponent(storage.NewR2Storage(common.PluginStorage)),
		appctx.WithComponent(tracer.NewJaeger(common.PluginTracer)),
		appctx.WithComponent(pubsub.NewNatsPubSub(common.PluginPubSub)),
		appctx.WithComponent(redisdb.NewRedisDB(common.PluginRedis, common.PluginRedis)),
		appctx.WithComponent(rpccaller.NewItemAPICaller(common.PluginItemAPI)),
		appctx.WithComponent(ginserver.NewServer(common.PluginGin)),
		appctx.WithComponent(grpcserver.NewServer(common.PluginGRPCServer)),
		appctx.WithComponent(grpcclient.NewClient(common.PluginGRPCClient, "")),
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

		common.AppStore = common.NewAppStore(appCtx)

		subscriber.StartPbEngine(appCtx)

		go startGRPCServer(appCtx)
		startGinServer(appCtx)
	},
}

func startGRPCServer(ac appctx.AppContext) {
	grpcServer := ac.MustGet(common.PluginGRPCServer).(core.GRPCServerComponent)
	server := grpcServer.GetServer()

	pb.RegisterUserLikeItemServiceServer(server, builder.BuildUserLikeItemGRPCService())

	grpcServer.Start(context.Background())
}

func startGinServer(ac appctx.AppContext) {
	ginServer := ac.MustGet(common.PluginGin).(core.GinComponent)
	router := ginServer.GetRouter()
	v1 := router.Group("/v1")
	setupRoutes(v1)

	ginServer.Start()
}

func setupRoutes(router *gin.RouterGroup) {
	userService := builder.BuildUserAPIService()
	itemService := builder.BuildItemAPIService()
	userLikeItemService := builder.BuildUserLikeItemAPIService()
	//userLikeItemRPCService := builder.BuildUserLikeItemRPCService()
	uploadService := builder.BuildUploadAPIService()
	authMiddleware := builder.BuildAuthMiddleware()

	router.Use(middleware.Recovery(common.AppStore.CTX))
	router.Static("/static", "./static")

	router.POST("/register", userService.Register())
	router.POST("/login", userService.Login())
	router.GET("/profile", authMiddleware, userService.Profile())

	uploads := router.Group("/upload", authMiddleware)
	{
		uploads.POST("", uploadService.Upload())
		uploads.POST("/local", uploadService.UploadLocal())
	}

	items := router.Group("/items", authMiddleware)
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

	//rpc := router.Group("/rpc")
	//{
	//	rpc.POST("/get_item_likes", userLikeItemRPCService.GetItemLikes())
	//}
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronUpdateItemLikedCountCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
