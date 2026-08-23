package entity

import "errors"

var (
	ErrParametrNoFound = errors.New("params not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrCodeInvalid     = errors.New("code invalid")
	ErrTokenInvalid    = errors.New("token invalid")
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenBlocked    = errors.New("token blocked")
)
