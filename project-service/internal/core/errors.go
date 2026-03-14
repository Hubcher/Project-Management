package core

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidProject  = errors.New("invalid project")
)
