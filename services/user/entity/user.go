package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
)

type UserRole int

const (
	RoleUser UserRole = 1 << iota
	RoleAdmin
	RoleMod
)

func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	case RoleMod:
		return "mod"
	default:
		return "user"
	}
}

func getRoleFromStr(value string) UserRole {
	switch value {
	case "admin":
		return RoleAdmin
	case "mod":
		return RoleMod
	default:
		return RoleUser
	}
}

func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New((fmt.Sprint("Failed to unmarshal JSONB value:", value)))
	}

	roleValue := string(bytes)
	*role = getRoleFromStr(roleValue)

	return nil
}

func (role *UserRole) Value() (driver.Value, error) {
	if role == nil {
		return nil, nil
	}

	return role.String(), nil
}

func (role *UserRole) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", role.String())), nil
}

type User struct {
	core.SQLModel
	Email     string      `json:"email" gorm:"column:email;"`
	Password  string      `json:"-" gorm:"column:password;"`
	Salt      string      `json:"-" gorm:"column:salt;"`
	LastName  string      `json:"last_name" gorm:"column:last_name;"`
	FirstName string      `json:"first_name" gorm:"column:first_name;"`
	Phone     string      `json:"phone" gorm:"column:phone;"`
	Avatar    *core.Image `json:"avatar" gorm:"column:avatar;"`
	Role      UserRole    `json:"role" gorm:"column:role;"`
	Status    int         `json:"status" gorm:"column:status;"`
}

func NewUser(firstName, lastName, email string) User {
	return User{
		SQLModel:  core.NewSQLModel(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     "",
		Avatar:    nil,
		Role:      RoleUser,
		Status:    0,
	}
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetUserID() int {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

func (u *User) Mask() {
	u.SQLModel.Mask(common.MaskTypeUser)
}
