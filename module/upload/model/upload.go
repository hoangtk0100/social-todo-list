package model

import (
	"github.com/hoangtk0100/social-todo-list/common"
)

type Upload struct {
	common.SQLModel `json:",inline"`
	common.Image    `json:",inline"`
}

func (Upload) TableName() string { return "uploads" }
