package psql

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/sausage-auth/internal/entity"
	"github.com/RenterRus/sausage-auth/internal/repo/psql/db"
)

// Confirmed implements db.Querier.
func (u *UserRepo) Confirmed(ctx context.Context, login *string) error {
	if login == nil || *login == "" {
		return fmt.Errorf("Confirmed: %w", entity.ErrParametrNoFound)
	}

	if err := u.Queries.Confirmed(ctx, login); err != nil {
		return fmt.Errorf("Confirmed.Confirmed: %w", err)
	}

	return nil
}

// Hash implements db.Querier.
func (u *UserRepo) OtpHash(ctx context.Context, login *string) (string, error) {
	if login == nil || *login == "" {
		return "", fmt.Errorf("Hash: %w", entity.ErrParametrNoFound)
	}

	hash, err := u.Queries.OtpHash(ctx, login)
	if err != nil {
		return "", fmt.Errorf("Hash.Hash: %w", err)
	}

	return hash, nil
}

// IsConfirmed implements db.Querier.
func (u *UserRepo) IsExist(ctx context.Context, login *string) (bool, error) {
	if login == nil || *login == "" {
		return false, fmt.Errorf("IsConfirmed: %w", entity.ErrParametrNoFound)
	}

	userLogin, err := u.Queries.IsExist(ctx, login)
	if err != nil {
		return false, fmt.Errorf("IsConfirmed.IsConfirmed: %w", err)
	}

	return userLogin, nil
}

// Register implements db.Querier.
func (u *UserRepo) Register(ctx context.Context, arg entity.RegisterParams) error {
	if arg.Login == "" {
		return fmt.Errorf("Register(login): %w", entity.ErrParametrNoFound)
	}

	if arg.Hash == "" {
		return fmt.Errorf("Register(hash): %w", entity.ErrParametrNoFound)
	}

	if arg.URL == "" {
		return fmt.Errorf("Register(url): %w", entity.ErrParametrNoFound)
	}

	err := u.Queries.Register(ctx, db.RegisterParams{
		Login: &arg.Login,
		Hash:  &arg.Hash,
		Url:   &arg.URL,
	})
	if err != nil {
		return fmt.Errorf("Register.Register: %w", err)
	}

	return nil
}

func (u *UserRepo) URLByUUID(ctx context.Context, uuid *string) (string, error) {
	if pointer.Get(uuid) == "" {
		return "", fmt.Errorf("LoginByUUID: %w", entity.ErrParametrNoFound)
	}

	login, err := u.Queries.URLByUUID(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("LoginByUUID.LoginByUUID: %w", err)
	}

	return login, nil
}
