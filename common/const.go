package common

import "time"

const (
	CurrentUser        = "current_user"
	PluginDBMain       = "mysql"
	PluginJWT          = "jwt"
	PluginR2           = "r2"
	PluginPubSub       = "pubsub"
	PluginItemAPI      = "item-api"
	PluginTracerJaeger = "social-todo-jaeger"
	PluginRedis        = "redis"
	PluginGin          = "gin"

	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"
)

type DBType int

const (
	DBTypeUser DBType = 1
	DBTypeItem DBType = 2
)

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}

func IsOwner(requester Requester, ownerId int) bool {
	return requester.GetUserId() == ownerId
}
