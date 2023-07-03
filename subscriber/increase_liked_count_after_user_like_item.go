package subscriber

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/services/item/repository/mysql"
)

type HasItemID interface {
	GetItemID() int
}

func IncreaseLikedCountAfterUserLikeItem(ac appctx.AppContext) core.SubJob {
	return core.SubJob{
		Name: "Increase liked count after user likes item",
		Hdl: func(ctx context.Context, msg *pubsub.Message) error {
			db := ac.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()
			data := msg.Data().(map[string]interface{})
			itemID := data["item_id"].(float64)

			if err := mysql.NewMySQLRepository(db).IncreaseLikedCount(ctx, int(itemID)); err != nil {
				return err
			}

			_ = msg.Ack()

			return nil
		},
	}
}
