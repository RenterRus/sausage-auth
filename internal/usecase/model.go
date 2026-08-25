package usecase

import (
	"context"

	"github.com/RenterRus/sausage-auth/internal/entity"
)

type RevokeType int

const (
	REVOKE_ONE RevokeType = iota
	REVOKE_ALL
)

type LoginRequest struct {
	Login     string
	UserAgent string
	Code      string
}

type RefreshRequest struct {
	Login        string
	UserAgent    string
	RefreshToken string
}

type Register interface {
	Registration(ctx context.Context, login string) (string, error)
	Confirmed(ctx context.Context, login, code string) error

	LoginOTP(ctx context.Context, req LoginRequest) (entity.Tokens, error)
	Refresh(ctx context.Context, req RefreshRequest) (entity.Tokens, error)

	Logout(ctx context.Context, sign *string, revokeType RevokeType) error

	// func(ctx, access) userID
	Validation(ctx context.Context, access string) (*string, error)
}
