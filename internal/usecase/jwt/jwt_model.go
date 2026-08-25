package jwt

import "github.com/RenterRus/sausage-auth/internal/entity"

type JWT interface {
	GenAccess(user_login string) (string, error)
	GenRefresh(user_login string) (string, error)

	Parse(token string) (entity.BaseJWT, error)
}
