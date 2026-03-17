package core

import "errors"

var (
	ErrAccountNotFound             = errors.New("account not found")
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrAccountInactive             = errors.New("account is inactive")
	ErrRefreshSessionNotFound      = errors.New("refresh session not found")
	ErrRefreshSessionAlreadyExists = errors.New("refresh session already exists")
	ErrRefreshTokenExpired         = errors.New("refresh token expired")
	ErrRefreshTokenRevoked         = errors.New("refresh token revoked")
)
