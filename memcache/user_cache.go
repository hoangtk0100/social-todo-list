package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/services/user/entity"
)

type RealUserRepository interface {
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

type userCache struct {
	repo     core.CacheComponent
	realRepo RealUserRepository
	once     *sync.Once
}

func NewUserCache(repo core.CacheComponent, realRepo RealUserRepository) *userCache {
	return &userCache{
		repo:     repo,
		realRepo: realRepo,
		once:     new(sync.Once),
	}
}

func (c *userCache) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	key := fmt.Sprintf("user-%d", id)

	var user entity.User
	err := c.repo.Get(ctx, key, &user)
	if err == nil && user.ID > 0 {
		return &user, nil
	}

	var userErr error
	c.once.Do(func() {
		realUser, err := c.realRepo.GetUserByID(ctx, id)
		if err != nil {
			log.Println(userErr)
			userErr = err
			return
		}

		user = *realUser
		_ = c.repo.Set(ctx, key, realUser, time.Hour*2)
	})

	if userErr != nil {
		return nil, userErr
	}

	err = c.repo.Get(ctx, key, &user)
	if err == nil && user.ID > 0 {
		return &user, nil
	}

	return nil, err
}
