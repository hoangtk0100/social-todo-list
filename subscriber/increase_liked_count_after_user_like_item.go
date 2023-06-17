package subscriber

import (
	"context"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/module/item/storage"
	"github.com/hoangtk0100/social-todo-list/pubsub"
	"gorm.io/gorm"
)

type HasItemId interface {
	GetItemId() int
}

func IncreaseLikedCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Name: "Increase liked count after user likes item",
		Hdl: func(ctx context.Context, msg *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			data := msg.Data().(map[string]interface{})
			itemId := data["item_id"].(float64)

			if err := storage.NewSQLStore(db).IncreaseLikedCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = msg.Ack()

			return nil
		},
	}
}
