package core

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUser       = errors.New("invalid user")
	ErrUserAlreadyExists = errors.New("user already exists")
)
