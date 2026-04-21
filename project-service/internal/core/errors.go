package core

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidProject  = errors.New("invalid project")
	ErrStageNotFound   = errors.New("project stage not found")
	ErrInvalidStage    = errors.New("invalid project stage")
	ErrMemberNotFound  = errors.New("project member not found")
	ErrInvalidMember   = errors.New("invalid project member")
	ErrEventNotFound   = errors.New("project event not found")
	ErrInvalidEvent    = errors.New("invalid project event")
	ErrAlreadyExists   = errors.New("resource already exists")
)
