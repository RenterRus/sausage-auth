package inmem

import (
	"context"
)

type AccessCacheRequest struct {
	Key   string
	Value string
}

type AccessCache interface {
	Set(ctx context.Context, req AccessCacheRequest) error
	Get(ctx context.Context, access string) (string, bool, error)
}
