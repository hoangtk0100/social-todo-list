package common

import (
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/cache"
	"github.com/hoangtk0100/app-context/core"
	"gorm.io/gorm"
)

var AppStore *appStore

type appStore struct {
	CTX        appctx.AppContext
	DB         *gorm.DB
	CacheDB    core.CacheComponent
	TokenMaker core.TokenMakerComponent
	Storage    core.StorageComponent
	PS         core.PubSubComponent

	ItemAPICaller ItemAPICaller
}

func NewAppStore(appCtx appctx.AppContext) *appStore {
	db := appCtx.MustGet(PluginDBMain).(core.GormDBComponent).GetDB()

	return &appStore{
		CTX:           appCtx,
		DB:            db,
		TokenMaker:    appCtx.MustGet(PluginTokenMaker).(core.TokenMakerComponent),
		PS:            appCtx.MustGet(PluginPubSub).(core.PubSubComponent),
		Storage:       appCtx.MustGet(PluginStorage).(core.StorageComponent),
		CacheDB:       cache.NewRedisCache(PluginRedis, appCtx),
		ItemAPICaller: appCtx.MustGet(PluginItemAPI).(ItemAPICaller),
	}
}
