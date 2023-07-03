package model

import (
	"github.com/hoangtk0100/app-context/core"
)

type Upload struct {
	core.SQLModel `json:",inline"`
	core.Image    `json:",inline"`
}

func (Upload) TableName() string { return "uploads" }
