package subscriber

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
)

func StartPbEngine(ac appctx.AppContext) {
	ps := ac.MustGet(common.PluginPubSub).(core.PubSubComponent)
	pbEngine := core.NewSubscribeEngine(common.PubSubEngineName, ps, ac)
	pbEngine.AddTopicJobs(
		common.TopicUserLikedItem,
		true,
		IncreaseLikedCountAfterUserLikeItem(ac),
	)

	pbEngine.AddTopicJobs(
		common.TopicUserUnlikedItem,
		true,
		DecreaseLikedCountAfterUserUnlikeItem(ac),
	)

	pbEngine.Start()
}
