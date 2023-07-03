package common

import (
	"time"
)

const (
	PluginDBMain       = "mysql"
	PluginJWT          = "jwt"
	PluginR2           = "r2"
	PluginPubSub       = "pubsub"
	PluginItemAPI      = "item-api"
	PluginTracerJaeger = "social-todo-jaeger"
	PluginRedis        = "redis"
	PluginGin          = "gin"

	PubSubEngineName = "pb-engine"

	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"

	HashPasswordFormat = "%s.%s"

	MaskTypeUser = 1
	MaskTypeItem = 2
)

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}
