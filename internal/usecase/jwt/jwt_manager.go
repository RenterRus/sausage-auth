package jwt

import (
	"fmt"
	"time"

	"github.com/RenterRus/sausage-profile/internal/entity"
	gojwt "github.com/golang-jwt/jwt/v5"
)

type jwtManager struct {
	access_exp  time.Duration
	refresh_exp time.Duration
	key         []byte
}

type JwtManagerConf struct {
	AccessExp  time.Duration
	RefreshExp time.Duration
	Key        []byte
}

func NewJWTManager(req JwtManagerConf) JWT {
	return &jwtManager{
		access_exp:  req.AccessExp,
		refresh_exp: req.RefreshExp,
		key:         req.Key,
	}
}

func (j *jwtManager) GenAccess(user_login string) (string, error) {
	token, err := j.generateToken(user_login, time.Now().Add(j.access_exp))
	if err != nil {
		return "", fmt.Errorf("GenAccess.generateToken: %w", err)
	}

	return token, nil
}

func (j *jwtManager) GenRefresh(user_login string) (string, error) {
	token, err := j.generateToken(user_login, time.Now().Add(j.refresh_exp))
	if err != nil {
		return "", fmt.Errorf("GenRefresh.generateToken: %w", err)
	}

	return token, nil
}

func (j *jwtManager) Parse(token string) (entity.BaseJWT, error) {
	tkn, err := j.parseToken(token)
	if err != nil {
		return entity.BaseJWT{}, fmt.Errorf("Parse.parseToken: %w", err)
	}

	sub, err := tkn.Claims.GetSubject()
	if err != nil {
		return entity.BaseJWT{}, fmt.Errorf("Parse.GetSubject: %w", err)
	}

	exp, err := tkn.Claims.GetExpirationTime()
	if err != nil {
		return entity.BaseJWT{}, fmt.Errorf("Parse.GetExpirationTime: %w", err)
	}

	return entity.BaseJWT{
		IsValid:   tkn.Valid,
		UserLogin: sub,
		IsExpired: exp.Before(time.Now()),
	}, nil
}

func (j *jwtManager) generateToken(userID string, expired time.Time) (string, error) {
	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS384, gojwt.MapClaims{
		"sub": userID,
		"exp": expired.Unix(),
	}).SignedString(j.key)
	if err != nil {
		return "", fmt.Errorf("generateToken: %w", err)
	}

	return token, nil
}

func (j *jwtManager) parseToken(tokenString string) (*gojwt.Token, error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parseToken: %w", err)
	}

	return token, nil
}
