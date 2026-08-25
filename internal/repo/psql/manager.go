package psql

import (
	"github.com/RenterRus/sausage-auth/internal/repo/psql/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	Queries *db.Queries
}

func NewDBManager(pgx *pgxpool.Pool) UsersRepo {
	return &UserRepo{
		Queries: db.New(pgx),
	}
}
