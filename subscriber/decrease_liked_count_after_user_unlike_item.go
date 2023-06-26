package subscriber

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
)

func DecreaseLikedCountAfterUserUnlikeItem(ac appctx.AppContext) subJob {
	return subJob{
		Name: "Decrease liked count after user unlikes item",
		Hdl: func(ctx context.Context, msg *pubsub.Message) error {
			db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
			data := msg.Data().(map[string]interface{})
			itemId := data["item_id"].(float64)

			if err := storage.NewSQLStore(db).DecreaseLikedCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = msg.Ack()

			return nil
		},
	}
}
