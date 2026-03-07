package core

import "errors"

var ErrBadArguments = errors.New("arguments are not acceptable")
var ErrAlreadyExists = errors.New("project already exists")
var ErrNotFound = errors.New("project is not found")
