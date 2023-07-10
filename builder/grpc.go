package builder

import (
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/pb"
	user_like_item_business "github.com/hoangtk0100/social-todo-list/services/userlikeitem/business"
	user_like_item_sql_repo "github.com/hoangtk0100/social-todo-list/services/userlikeitem/repository/mysql"
	user_like_tem_grpc "github.com/hoangtk0100/social-todo-list/services/userlikeitem/transport/grpc"
)

func BuildUserLikeItemGRPCService() pb.UserLikeItemServiceServer {
	repo := user_like_item_sql_repo.NewMySQLRepository(common.AppStore.DB)
	business := user_like_item_business.NewBusiness(repo, common.AppStore.PS)
	service := user_like_tem_grpc.NewService(business)
	return service
}

func buildUserLikeItemGRPCClient() pb.UserLikeItemServiceClient {
	client := common.AppStore.CTX.MustGet(common.PluginGRPCClient).(core.GRPCClientComponent)
	conn := client.Dial()
	return pb.NewUserLikeItemServiceClient(conn)
}
