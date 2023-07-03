package cmd

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/datastore/gormdb"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/spf13/cobra"
)

var cronUpdateItemLikedCountCmd = &cobra.Command{
	Use:   "update-count",
	Short: "Update item liked count for all todo_items records",
	Run: func(cmd *cobra.Command, args []string) {
		service := appctx.NewAppContext(
			appctx.WithName("social-todo-list"),
			appctx.WithComponent(gormdb.NewGormDB(common.PluginDBMain, common.PluginDBMain)),
		)

		log := service.Logger("update-count-service")

		if err := service.Load(); err != nil {
			log.Fatal(err)
		}

		db := service.MustGet(common.PluginDBMain).(core.GormDBComponent).GetDB()

		query := `UPDATE todo_items ti
		INNER JOIN (
			SELECT item_id, COUNT(item_id) AS count
			FROM user_like_items
			GROUP BY item_id
		) c ON c.item_id = ti.id
		SET ti.liked_count = c.count`

		db.Exec(query)
		log.Info("Updated item liked count")
	},
}
