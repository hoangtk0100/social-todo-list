package model

import (
	"time"

	"github.com/hoangtk0100/app-context/core"
)

type Like struct {
	UserID    int              `json:"user_id" gorm:"column:user_id;"`
	ItemID    int              `json:"item_id" gorm:"column:item_id;"`
	CreatedAt *time.Time       `json:"created_at" gorm:"column:created_at;"`
	User      *core.SimpleUser `json:"-" gorm:"foreignKey:UserID;"`
}

func (Like) TableName() string { return "user_like_items" }

func (l *Like) GetItemID() int { return l.ItemID }
func (l *Like) GetUserID() int { return l.UserID }
