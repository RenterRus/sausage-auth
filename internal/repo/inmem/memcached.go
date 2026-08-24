package inmem

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type accessCache struct {
	accessExp time.Duration
	client    *memcache.Client
}

type AccessCacheConf struct {
	AccessExp time.Duration
	Client    *memcache.Client
}

func NewAccessCache(conf AccessCacheConf) AccessCache {
	return &accessCache{
		client:    conf.Client,
		accessExp: conf.AccessExp,
	}
}

func (a *accessCache) Set(ctx context.Context, req AccessCacheRequest) error {
	if err := a.client.Add(&memcache.Item{
		Key:        req.Access,
		Value:      []byte(req.UserUUID),
		Expiration: int32(a.accessExp),
	}); err != nil {
		return fmt.Errorf("Cahce.Set: %w", err)
	}

	return nil
}

func (a *accessCache) Get(ctx context.Context, access string) (string, bool, error) {
	acc, err := a.client.Get(access)
	if err != nil {
		return "", false, fmt.Errorf("Cache.Get: %w", err)
	}

	return string(acc.Value), true, nil
}
