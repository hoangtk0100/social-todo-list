package model

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/hoangtk0100/social-todo-list/common"
)

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted",
	)
)

const EntityName = "User"

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
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"-" gorm:"column:password;"`
	Salt      string   `json:"-" gorm:"column:salt;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	Phone     string   `json:"phone" gorm:"column:phone;"`
	Role      UserRole `json:"role" gorm:"column:role;"`
	Status    int      `json:"status" gorm:"column:status;"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

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
