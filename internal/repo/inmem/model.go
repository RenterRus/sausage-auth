package inmem

import (
	"context"
)

type AccessCacheRequest struct {
	Access   string
	UserUUID string
}

type AccessCache interface {
	Set(ctx context.Context, req AccessCacheRequest) error
	Get(ctx context.Context, access string) (string, bool, error)
}
