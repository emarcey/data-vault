package dependencies

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"emarcey/data-vault/common"
	"emarcey/data-vault/database"
)

type AccessTokenCacheUpdate struct {
	id          string
	accessToken *common.AccessToken
	updateType  common.CacheUpdateType
}

type AccessTokenCache struct {
	m            sync.Mutex
	logger       *logrus.Logger
	accessTokens map[string]*common.AccessToken
	updates      chan AccessTokenCacheUpdate
}

func (u *AccessTokenCache) Get(id string) *common.AccessToken {
	accessToken, ok := u.accessTokens[id]
	if !ok {
		return nil
	}
	return accessToken
}

func (u *AccessTokenCache) Delete(id string) {
	u.updates <- AccessTokenCacheUpdate{
		id:          id,
		accessToken: nil,
		updateType:  common.CACHE_DELETE,
	}
}

func (u *AccessTokenCache) Add(id string, accessToken *common.AccessToken) {
	u.updates <- AccessTokenCacheUpdate{
		id:          id,
		accessToken: accessToken,
		updateType:  common.CACHE_ADD,
	}
}

func (u *AccessTokenCache) handleUpdate(msg AccessTokenCacheUpdate) {
	u.m.Lock()
	defer u.m.Unlock()

	switch msg.updateType {
	case common.CACHE_ADD:
		u.accessTokens[msg.id] = msg.accessToken
	case common.CACHE_DELETE:
		delete(u.accessTokens, msg.id)
	default:
		u.logger.Errorf("Unexpected message in accessToken cache: %+v", msg)
	}
}

func (u *AccessTokenCache) ProcessUpdates(ctx context.Context) {
	for true {
		select {
		case <-ctx.Done():
			u.logger.Debug("Context canceled. Closing AccessTokenCache")
			close(u.updates)
			return
		case msg := <-u.updates:
			u.handleUpdate(msg)
		}
	}
}

func (u *AccessTokenCache) handleRefresh(ctx context.Context, db *database.DatabaseEngine) error {
	u.m.Lock()
	defer u.m.Unlock()
	authAccessTokens, err := database.SelectAccessTokensForAuth(ctx, db)
	if err != nil {
		return err
	}
	u.accessTokens = authAccessTokens
	return nil
}

func (u *AccessTokenCache) Refresh(ctx context.Context, db *database.DatabaseEngine, dataRefreshSeconds int) {
	timer := time.NewTicker(time.Duration(dataRefreshSeconds) * time.Second)
	for true {
		select {
		case <-ctx.Done():
			u.logger.Debug("Context canceled. Closing AccessTokenCache")
			close(u.updates)
			return
		case <-timer.C:
			err := u.handleRefresh(ctx, db)
			if err != nil {
				u.logger.Errorf("Error in SelectAccessTokensForAuth refresh: %v", err)
			}
		}
	}
}

func NewAccessTokenCache(ctx context.Context, logger *logrus.Logger, db *database.DatabaseEngine, dataRefreshSeconds int) (*AccessTokenCache, error) {
	var m sync.Mutex
	accessTokenCache := &AccessTokenCache{
		m:            m,
		logger:       logger,
		accessTokens: make(map[string]*common.AccessToken),
		updates:      make(chan AccessTokenCacheUpdate, 10),
	}

	err := accessTokenCache.handleRefresh(ctx, db)
	if err != nil {
		return nil, err
	}

	go accessTokenCache.ProcessUpdates(ctx)
	go accessTokenCache.Refresh(ctx, db, dataRefreshSeconds)

	return accessTokenCache, nil
}
