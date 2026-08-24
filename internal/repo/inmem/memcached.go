package inmem

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	hash := sha256.Sum256([]byte(req.Access))

	if err := a.client.Add(&memcache.Item{
		Key:        hex.EncodeToString(hash[:]),
		Value:      []byte(req.UserUUID),
		Expiration: int32(time.Now().Add(a.accessExp).Unix()),
	}); err != nil {
		return fmt.Errorf("Cahce.Set: %w", err)
	}

	return nil
}

func (a *accessCache) Get(ctx context.Context, access string) (string, bool, error) {
	hash := sha256.Sum256([]byte(access))

	acc, err := a.client.Get(hex.EncodeToString(hash[:]))
	if err != nil {
		return "", false, fmt.Errorf("Cache.Get: %w", err)
	}

	return string(acc.Value), true, nil
}
