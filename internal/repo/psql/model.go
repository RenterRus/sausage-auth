package psql

import (
	"context"

	"github.com/RenterRus/sausage-auth/internal/entity"
)

type GetRefreshReq struct {
	Login     *string
	Hash      *string
	UserAgent *string
}
type UsersRepo interface {
	// register
	Confirmed(ctx context.Context, login *string) error
	OtpHash(ctx context.Context, login *string) (string, error)
	IsExist(ctx context.Context, login *string) (bool, error)
	Register(ctx context.Context, arg entity.RegisterParams) error

	// session
	GetRefreshToken(ctx context.Context, req GetRefreshReq) (entity.GetRefreshTokenRow, error)
	SetBlockRefresh(ctx context.Context, refreshHash *string) error
	SetRefreshHash(ctx context.Context, arg entity.SetRefreshHashParams) error

	GetUUIDByLogin(ctx context.Context, login *string) (string, error)
	URLByUUID(ctx context.Context, uuid *string) (string, error)
	RemoveRefreshByHash(ctx context.Context, refreshHash *string) error
	RemoveRefreshByLogin(ctx context.Context, userLogin *string) ([]string, error)

	DeleteOldRefresh(ctx context.Context, arg entity.DeleteOldRefreshParams) ([]string, error)

	UpdateLastSighUp(ctx context.Context, userLogin *string) error
}
