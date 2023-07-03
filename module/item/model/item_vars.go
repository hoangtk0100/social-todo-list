package model

import (
	"strings"

	"github.com/hoangtk0100/app-context/core"
)

type TodoItemCreation struct {
	ID          int          `json:"id" gorm:"column:id;"`
	UserID      int          `json:"-" gorm:"column:user_id;"`
	Title       string       `json:"title" gorm:"column:title;"`
	Description string       `json:"description" gorm:"column:description;"`
	Image       *core.Images `json:"image" gorm:"column:image"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)

	if i.Title == "" {
		return ErrTitleEmpty
	}

	return nil
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string      `json:"title" gorm:"column:title;"`
	Description *string      `json:"description" gorm:"column:description;"`
	Status      *string      `json:"status" gorm:"column:status;"`
	Image       *core.Images `json:"image" gorm:"column:image"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }
