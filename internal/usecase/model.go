package usecase

import (
	"context"

	"github.com/RenterRus/sausage-profile/internal/entity"
)

type RevokeType int

const (
	REVOKE_ONE RevokeType = iota
	REVOKE_ALL
)

type Register interface {
	Registration(ctx context.Context, login string) (string, error)
	Confirmed(ctx context.Context, login, code string) error

	LoginOTP(ctx context.Context, login, userAgent, code string) (entity.Tokens, error)
	Refresh(ctx context.Context, login, userAgent, refreshToken string) (entity.Tokens, error)

	Logout(ctx context.Context, login, hash *string, revokeType RevokeType) error

	// func(ctx, access) userID
	Validation(ctx context.Context, access string) (*string, error)
}
