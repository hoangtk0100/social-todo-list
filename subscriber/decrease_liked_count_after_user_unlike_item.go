package subscriber

import (
	"context"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
	"github.com/hoangtk0100/social-todo-list/pubsub"
	"gorm.io/gorm"
)

func DecreaseLikedCountAfterUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Name: "Decrease liked count after user unlikes item",
		Hdl: func(ctx context.Context, msg *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
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
