package cache

import (
	"context"
	"time"

	goservice "github.com/200Lab-Education/go-sdk"
	rdcache "github.com/go-redis/cache/v9"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	store *rdcache.Cache
}

func NewRedisCache(sc goservice.ServiceContext) *redisCache {
	client := sc.MustGet(common.PluginRedis).(*redis.Client)

	c := rdcache.New(&rdcache.Options{
		Redis:      client,
		LocalCache: rdcache.NewTinyLFU(1000, time.Minute),
	})

	return &redisCache{store: c}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&rdcache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}
