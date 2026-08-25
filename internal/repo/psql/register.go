package psql

import (
	"context"
	"fmt"

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
	if arg.Login == nil || *arg.Login == "" {
		return fmt.Errorf("Register(login): %w", entity.ErrParametrNoFound)
	}

	if arg.Hash == nil || *arg.Hash == "" {
		return fmt.Errorf("Register(hash): %w", entity.ErrParametrNoFound)
	}

	if arg.Link == nil || *arg.Link == "" {
		return fmt.Errorf("Register(link): %w", entity.ErrParametrNoFound)
	}

	err := u.Queries.Register(ctx, db.RegisterParams{
		Login: arg.Login,
		Hash:  arg.Hash,
		Link:  arg.Link,
	})
	if err != nil {
		return fmt.Errorf("Register.Register: %w", err)
	}

	return nil
}
