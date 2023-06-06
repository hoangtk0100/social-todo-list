package common

import "log"

const (
	CurrentUser = "current_user"
)

func Recovery() {
	if r := recover(); r != nil {
		log.Println("Recovered:", r)
	}
}

type TokenPayload struct {
	UId   int    `json:"user_id"`
	URole string `json:"role"`
}

func (p TokenPayload) UserId() int {
	return p.UId
}

func (p TokenPayload) Role() string {
	return p.URole
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}

func IsOwner(requester Requester, ownerId int) bool {
	return requester.GetUserId() == ownerId
}
