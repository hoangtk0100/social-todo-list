package subscriber

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
)

func StartPbEngine(ac appctx.AppContext) {
	pbEngine := core.NewSubscriberEngine(common.PubSubEngineName, common.AppStore.PS, ac)
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
