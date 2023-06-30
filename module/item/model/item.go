package model

import (
	"github.com/hoangtk0100/social-todo-list/common"
)

type TodoItem struct {
	common.SQLModel
	UserId      int                `json:"-" gorm:"column:user_id;"`
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      string             `json:"status" gorm:"column:status;"`
	LikedCount  int                `json:"liked_count" gorm:"-"`
	Image       *common.Images     `json:"image" gorm:"column:image"`
	Owner       *common.SimpleUser `json:"owner" gorm:"foreignKey:UserId;"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (i *TodoItem) Mask() {
	i.SQLModel.Mask(common.DBTypeItem)
	if value := i.Owner; value != nil {
		value.Mask()
	}
}
