package builder

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/memcache"
	"github.com/hoangtk0100/social-todo-list/middleware"
	item_business "github.com/hoangtk0100/social-todo-list/services/item/business"
	item_grpc_repo "github.com/hoangtk0100/social-todo-list/services/item/repository/grpc"
	item_sql_repo "github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
	item_api "github.com/hoangtk0100/social-todo-list/services/item/transport/api"
	upload_business "github.com/hoangtk0100/social-todo-list/services/upload/business"
	upload_sql_repo "github.com/hoangtk0100/social-todo-list/services/upload/repository/mysql"
	upload_api "github.com/hoangtk0100/social-todo-list/services/upload/transport/api"
	user_business "github.com/hoangtk0100/social-todo-list/services/user/business"
	user_sql_repo "github.com/hoangtk0100/social-todo-list/services/user/repository/mysql"
	user_api "github.com/hoangtk0100/social-todo-list/services/user/transport/api"
	user_like_item_business "github.com/hoangtk0100/social-todo-list/services/userlikeitem/business"
	user_like_item_sql_repo "github.com/hoangtk0100/social-todo-list/services/userlikeitem/repository/mysql"
	//user_like_item_grpc_repo "github.com/hoangtk0100/social-todo-list/services/userlikeitem/repository/grpc"
	user_like_item_api "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/api"
	user_like_item_rpc "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/rpc"
)

type (
	UserService interface {
		Login() gin.HandlerFunc
		Register() gin.HandlerFunc
		Profile() gin.HandlerFunc
	}

	ItemService interface {
		CreateItem() gin.HandlerFunc
		DeleteItem() gin.HandlerFunc
		GetItem() gin.HandlerFunc
		ListItem() gin.HandlerFunc
		UpdateItem() gin.HandlerFunc
	}

	UserLikeItemService interface {
		ListLikedUsers() gin.HandlerFunc
		LikeItem() gin.HandlerFunc
		UnlikeItem() gin.HandlerFunc
	}

	UserLikeItemRPCService interface {
		GetItemLikes() gin.HandlerFunc
	}

	UploadService interface {
		Upload() gin.HandlerFunc
		UploadLocal() gin.HandlerFunc
	}
)

func BuildAuthMiddleware() gin.HandlerFunc {
	repo := user_sql_repo.NewMySQLRepository(common.AppStore.DB)
	authCache := memcache.NewUserCache(common.AppStore.CacheDB, repo)
	return middleware.RequireAuth(authCache, common.AppStore.TokenMaker)
}

func BuildUserAPIService() UserService {
	repo := user_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := user_business.NewBusiness(repo, common.AppStore.TokenMaker)
	userService := user_api.NewService(common.AppStore.CTX, business)
	return userService
}

func BuildItemAPIService() ItemService {
	// itemAPICaller := common.AppStore.ItemAPICaller
	// likeRepo := rpc.NewItemAPIClient(itemAPICaller.GetServiceURL(), common.AppStore.CTX.Logger("rpc.itemlikes"))
	likeClient := item_grpc_repo.NewClient(buildUserLikeItemGRPCClient())
	repo := item_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := item_business.NewBusiness(repo, likeClient)
	service := item_api.NewService(common.AppStore.CTX, business)
	return service
}

func BuildUserLikeItemAPIService() UserLikeItemService {
	repo := user_like_item_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := user_like_item_business.NewBusiness(repo, common.AppStore.PS)
	service := user_like_item_api.NewService(common.AppStore.CTX, business)
	return service
}

func BuildUserLikeItemRPCService() UserLikeItemRPCService {
	repo := user_like_item_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := user_like_item_business.NewBusiness(repo, common.AppStore.PS)
	service := user_like_item_rpc.NewService(common.AppStore.CTX, business)
	return service
}

func BuildUploadAPIService() UploadService {
	repo := upload_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := upload_business.NewBusiness(repo, common.AppStore.Storage)
	service := upload_api.NewService(common.AppStore.CTX, business)
	return service
}
