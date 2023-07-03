package common

import (
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
)

func GetRequesterID(requester core.Requester) int {
	uid, _ := util.UIDFromString(requester.GetUID())
	return int(uid.GetLocalID())
}
