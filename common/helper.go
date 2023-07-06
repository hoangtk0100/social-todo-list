package common

import (
	"time"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
)

type ItemAPICaller interface {
	GetServiceURL() string
}

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func GetRequesterID(requester core.Requester) int {
	uid, _ := util.UIDFromString(requester.GetUID())
	return int(uid.GetLocalID())
}
