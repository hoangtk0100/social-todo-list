package common

const (
	PluginDBMain     = "mysql"
	PluginTokenMaker = "jwt"
	PluginStorage    = "r2"
	PluginPubSub     = "pubsub"
	PluginItemAPI    = "item-api"
	PluginTracer     = "social-todo-jaeger"
	PluginRedis      = "redis"
	PluginGin        = "gin"

	PubSubEngineName = "pb-engine"

	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"

	HashPasswordFormat = "%s.%s"

	MaskTypeUser = 1
	MaskTypeItem = 2
)
