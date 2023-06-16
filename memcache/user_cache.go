package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hoangtk0100/social-todo-list/module/user/model"
	"github.com/hoangtk0100/social-todo-list/plugin/cache"
)

type RealUserStore interface {
	FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type userCache struct {
	store     cache.Cache
	realStore RealUserStore
	once      *sync.Once
}

func NewUserCache(store cache.Cache, realStore RealUserStore) *userCache {
	return &userCache{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (c *userCache) FindUser(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.User, error) {
	userId := conds["id"].(int)
	key := fmt.Sprintf("user-%d", userId)

	var user model.User
	err := c.store.Get(ctx, key, &user)
	if err == nil && user.Id > 0 {
		return &user, nil
	}

	var userErr error
	c.once.Do(func() {
		realUser, err := c.realStore.FindUser(ctx, conds, moreInfo...)
		if err != nil {
			log.Println(userErr)
			userErr = err
			return
		}

		user = *realUser
		_ = c.store.Set(ctx, key, realUser, time.Hour*2)
	})

	if userErr != nil {
		return nil, userErr
	}

	err = c.store.Get(ctx, key, &user)
	if err == nil && user.Id > 0 {
		return &user, nil
	}

	return nil, err
}
