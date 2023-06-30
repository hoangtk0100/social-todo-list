package model

import "github.com/hoangtk0100/social-todo-list/common"

type UserCreate struct {
	common.SQLModel `json:",inline"`
	FirstName       string `json:"first_name" gorm:"column:first_name;"`
	LastName        string `json:"last_name" gorm:"column:last_name;"`
	Email           string `json:"email" gorm:"column:email;"`
	Password        string `json:"password" gorm:"column:password;"`
	Role            string `json:"-" gorm:"column:role;"`
	Salt            string `json:"-" gorm:"column:salt;"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}
