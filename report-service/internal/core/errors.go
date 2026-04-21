package core

import "errors"

var (
	ErrReportNotFound = errors.New("daily report not found")
	ErrInvalidReport  = errors.New("invalid daily report")
	ErrEntryNotFound  = errors.New("daily report entry not found")
	ErrInvalidEntry   = errors.New("invalid daily report entry")
	ErrCommentNotFound = errors.New("daily report comment not found")
	ErrInvalidComment  = errors.New("invalid daily report comment")
	ErrAlreadyExists   = errors.New("resource already exists")
)