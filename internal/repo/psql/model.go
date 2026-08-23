package psql

import (
	"context"

	"github.com/RenterRus/sausage-profile/internal/entity"
)

type UsersRepo interface {
	// register
	Confirmed(ctx context.Context, login *string) error
	Hash(ctx context.Context, login *string) (string, error)
	IsExist(ctx context.Context, login *string) (bool, error)
	Register(ctx context.Context, arg entity.RegisterParams) error

	// session
	GetRefreshToken(ctx context.Context, login, hash *string) (entity.GetRefreshTokenRow, error)
	SetBlockRefresh(ctx context.Context, refreshHash *string) error
	SetRefreshHash(ctx context.Context, arg entity.SetRefreshHashParams) error

	GetUUIDByLogin(ctx context.Context, login *string) (string, error)
	RemoveRefreshByHash(ctx context.Context, refreshHash *string) error
	RemoveRefreshByLogin(ctx context.Context, userLogin *string) ([]string, error)

	DeleteOldRefresh(ctx context.Context, arg entity.DeleteOldRefreshParams) ([]string, error)

	UpdateLastSighUp(ctx context.Context, userLogin *string) error
}
