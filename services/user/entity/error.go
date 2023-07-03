package entity

import "errors"

var (
	ErrEmailOrPasswordInvalid = errors.New("email or password invalid")
	ErrEmailExisted           = errors.New("email has already existed")
	ErrCannotGetUser          = errors.New("cannot get user info")
	ErrCannotGetUsers         = errors.New("cannot get users info")
	ErrCannotCreateUser       = errors.New("cannot create new user")
)
