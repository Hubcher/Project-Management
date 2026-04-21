package core

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExists      = errors.New("account already exists")
	ErrAccountNotFound    = errors.New("account not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInactiveAccount    = errors.New("inactive account")
	ErrInvalidArgument    = errors.New("invalid argument")
)
