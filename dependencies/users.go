package dependencies

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/emarcey/data-vault/common"
	"github.com/emarcey/data-vault/database"
)

type UserCacheUpdate struct {
	id         string
	user       *common.User
	updateType common.CacheUpdateType
}

type UserCache struct {
	m       sync.Mutex
	logger  *logrus.Logger
	users   map[string]*common.User
	updates chan UserCacheUpdate
}

func (u *UserCache) Get(id string) *common.User {
	user, ok := u.users[id]
	if !ok {
		return nil
	}
	return user
}

func (u *UserCache) Delete(id string) {
	u.updates <- UserCacheUpdate{
		id:         id,
		user:       nil,
		updateType: common.CACHE_DELETE,
	}
}

func (u *UserCache) Add(id string, user *common.User) {
	u.updates <- UserCacheUpdate{
		id:         id,
		user:       user,
		updateType: common.CACHE_ADD,
	}
}

func (u *UserCache) handleUpdate(msg UserCacheUpdate) {
	u.m.Lock()
	defer u.m.Unlock()

	switch msg.updateType {
	case common.CACHE_ADD:
		u.users[msg.id] = msg.user
	case common.CACHE_DELETE:
		delete(u.users, msg.id)
	default:
		u.logger.Errorf("Unexpected message in user cache: %+v", msg)
	}
}

func (u *UserCache) ProcessUpdates(ctx context.Context) {
	for true {
		select {
		case <-ctx.Done():
			u.logger.Debug("Context canceled. Closing UserCache")
			close(u.updates)
			return
		case msg := <-u.updates:
			u.handleUpdate(msg)
		}
	}
}

func (u *UserCache) handleRefresh(ctx context.Context, db *database.DatabaseEngine) error {
	u.m.Lock()
	defer u.m.Unlock()
	authUsers, err := database.SelectUsersForAuth(ctx, db)
	if err != nil {
		return err
	}
	u.users = authUsers
	return nil
}

func (u *UserCache) Refresh(ctx context.Context, db *database.DatabaseEngine, dataRefreshSeconds int) {
	timer := time.NewTicker(time.Duration(dataRefreshSeconds) * time.Second)
	for true {
		select {
		case <-ctx.Done():
			u.logger.Debug("Context canceled. Closing UserCache")
			close(u.updates)
			return
		case <-timer.C:
			err := u.handleRefresh(ctx, db)
			if err != nil {
				u.logger.Errorf("Error in SelectUsersForAuth refresh: %v", err)
			}
		}
	}
}

func NewUserCache(ctx context.Context, logger *logrus.Logger, db *database.DatabaseEngine, dataRefreshSeconds int) (*UserCache, error) {
	var m sync.Mutex
	userCache := &UserCache{
		m:       m,
		logger:  logger,
		users:   make(map[string]*common.User),
		updates: make(chan UserCacheUpdate, 10),
	}

	err := userCache.handleRefresh(ctx, db)
	if err != nil {
		return nil, err
	}

	go userCache.ProcessUpdates(ctx)
	go userCache.Refresh(ctx, db, dataRefreshSeconds)

	return userCache, nil
}

func NewMockUserCache(logger *logrus.Logger, users map[string]*common.User) *UserCache {
	var m sync.Mutex
	return &UserCache{
		m:       m,
		logger:  logger,
		users:   users,
		updates: make(chan UserCacheUpdate, 10),
	}
}
