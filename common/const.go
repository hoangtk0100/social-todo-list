package common

const (
	PluginDBMain     = "mysql"
	PluginTokenMaker = "paseto"
	PluginStorage    = "r2"
	PluginPubSub     = "pubsub"
	PluginItemAPI    = "item-api"
	PluginTracer     = "social-todo-jaeger"
	PluginRedis      = "redis"
	PluginGin        = "gin"
	PluginGRPCServer = "grpc-server"
	PluginGRPCClient = "grpc-client"

	PubSubEngineName = "pb-engine"

	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"

	HashPasswordFormat = "%s.%s"

	MaskTypeUser = 1
	MaskTypeItem = 2
)
