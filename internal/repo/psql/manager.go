package psql

import (
	"context"
	"fmt"

	"github.com/RenterRus/sausage-auth/internal/repo/psql/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	Queries *db.Queries
}

func NewDBManager(conn string) (UsersRepo, error) {
	pgx, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, fmt.Errorf("NewDBManager.New: %w", err)
	}

	return &UserRepo{
		Queries: db.New(pgx),
	}, nil
}
