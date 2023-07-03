package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hoangtk0100/app-context/component/cache"
	"github.com/hoangtk0100/social-todo-list/module/user/model"
)

type RealUserStore interface {
	GetUserByID(ctx context.Context, id int) (*model.User, error)
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

func (c *userCache) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	key := fmt.Sprintf("user-%d", id)

	var user model.User
	err := c.store.Get(ctx, key, &user)
	if err == nil && user.ID > 0 {
		return &user, nil
	}

	var userErr error
	c.once.Do(func() {
		realUser, err := c.realStore.GetUserByID(ctx, id)
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
	if err == nil && user.ID > 0 {
		return &user, nil
	}

	return nil, err
}
