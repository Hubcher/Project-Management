package core

import "errors"

var (
	ErrUserExists         = errors.New("User already exists")
	ErrUserNotFound       = errors.New("User not found")
	ErrAppNotFound        = errors.New("App not found")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrInvalidAppID       = errors.New("Invalid App ID")
)
