package cmd

import (
	"log"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/plugin/datastore/gormdb"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var cronUpdateItemLikedCountCmd = &cobra.Command{
	Use:   "update-count",
	Short: "Update item liked count for all todo_items records",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(
			goservice.WithName("social-todo-list"),
			goservice.WithVersion("1.0.0"),
			goservice.WithInitRunnable(gormdb.NewGormDB("main.mysql", common.PluginDBMain)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		db := service.MustGet(common.PluginDBMain).(*gorm.DB)

		query := `UPDATE todo_items ti
		INNER JOIN (
			SELECT item_id, COUNT(item_id) AS count
			FROM user_like_items
			GROUP BY item_id
		) c ON c.item_id = ti.id
		SET ti.liked_count = c.count`

		db.Exec(query)
		log.Println("Updated item liked count")
	},
}
